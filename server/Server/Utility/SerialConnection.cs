using System.Collections.Concurrent;
using System.IO.Ports;

namespace Server.Utility;

public abstract class SerialConnection : IDisposable
{
    private const ulong KILL_CONNECTION_AFTER_N_EMPTY_POLLS = 120;

    private readonly ReaderWriterLockSlim sendQueueLock;
    private readonly Queue<string> sendQueue;

    private Thread? communicationThread;

    private bool disposedValue;
    private bool isAlive;

    private string? activeSerialPortName;

    protected readonly ILogger<SerialConnection> logger;
    protected readonly ConcurrentCircularBuffer<string> concurrentCircularBuffer;

    /// <summary>
    /// Create a new serial connection
    /// </summary>
    /// <param name="logger">Logger</param>
    public SerialConnection(ILogger<SerialConnection> logger)
    {
        this.logger = logger;
        sendQueue = new Queue<string>();
        sendQueueLock = new ReaderWriterLockSlim();
        concurrentCircularBuffer = new ConcurrentCircularBuffer<string>(5_000);
    }

    /// <summary>
    /// Connect to the serial port
    /// </summary>
    /// <param name="portName">Name of the port to connect to</param>
    /// <returns>True if a connection was established successfully, false when not</returns>
    public bool Connect(string portName)
    {
        if (activeSerialPortName != null)
        {
            logger.LogError("Cannot connect because no serial port name was specified");
            return false;
        }

        activeSerialPortName = portName;

        communicationThread = new Thread(Run);
        communicationThread.Start();

        return communicationThread.IsAlive;
    }

    /// <summary>
    /// Disconnect from the serial port
    /// </summary>
    /// <returns>True if the application successfully disconnected from the serial port, false when not</returns>
    public bool Disconnect()
    {
        if (activeSerialPortName == null || !isAlive)
        {
            return true;
        }

        isAlive = false;
        communicationThread?.Join();
        sendQueueLock.EnterWriteLock();
        sendQueue.Clear();
        sendQueueLock.ExitWriteLock();
        return communicationThread?.IsAlive == false;
    }

    /// <summary>
    /// Dump all data in the buffer
    /// </summary>
    /// <returns>Raw data dump of the buffer's contents</returns>
    public List<string?> Dump() => concurrentCircularBuffer.ReadAll();

    /// <summary>
    /// Enqueue a message to be sent on the communication thread
    /// </summary>
    /// <param name="message">Message to send</param>
    public void SendRaw(string message)
    {
        sendQueueLock.EnterWriteLock();
        sendQueue.Enqueue(message);
        sendQueueLock.ExitWriteLock();
    }

    /// <summary>
    /// Dispose pattern implementation
    /// </summary>
    public void Dispose()
    {
        Dispose(disposing: true);
        GC.SuppressFinalize(this);
    }

    /// <summary>
    /// Gracefully dispose of the communication thread
    /// </summary>
    /// <param name="disposing">Dispose pattern implementation</param>
    protected virtual void Dispose(bool disposing)
    {
        if (!disposedValue)
        {
            if (disposing)
            {
                isAlive = false;
                communicationThread?.Join();
            }

            disposedValue = true;
        }
    }

    /// <summary>
    /// Communication thread loop
    /// </summary>
    private void Run()
    {
        if (activeSerialPortName == null)
        {
            logger.LogError("Cannot start thread because no serial port name was specified");
            return;
        }

        var port = new SerialPort(activeSerialPortName)
        {
            ReadTimeout = 500,
            WriteTimeout = 500
        };

        isAlive = true;

        // Unable to open the port - might be in use or just not available in general
        try
        {
            logger.LogInformation("Attemping to open serial connection");
            port.Open();
            logger.LogInformation("Serial connection established");
        }
        catch (UnauthorizedAccessException)
        {
            logger.LogCritical("Unable to connect to serial port - unauthorized");

            // Ensure the thread kills itself immediately
            isAlive = false;
        }

        ulong pollCount = 0;

        while (isAlive)
        {
            pollCount++;
            var hasActionOccurred = false;

            // Attempt to write an outgoing message
            sendQueueLock.EnterReadLock();

            if (sendQueue.TryDequeue(out var data))
            {
                if (data != null)
                {
                    logger.LogDebug("Writing to port");

                    port.Write(data);
                    hasActionOccurred = true;
                }
            }

            sendQueueLock.ExitReadLock();

            // Attempt to read incoming messages
            try
            {
                concurrentCircularBuffer.Write(port.ReadLine());
                hasActionOccurred = true;
            }
            catch (TimeoutException) { }

            // Ensure the thread will exist automatically after a period of inactivity
            if (hasActionOccurred)
            {
                pollCount = 0;
            }
            else if (pollCount > KILL_CONNECTION_AFTER_N_EMPTY_POLLS)
            {
                logger.LogWarning("Serial communication thread will kill itself because of inactivity");
                isAlive = false;
            }
        }

        port.Close();
        activeSerialPortName = null;

        logger.LogInformation("Serial port has been closed and the thread has been stopped");
    }
}
