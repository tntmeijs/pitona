using System.IO.Ports;

namespace server.Utility;

public class SerialConnection : IDisposable
{
    private const ulong KILL_CONNECTION_AFTER_N_EMPTY_POLLS = 120;

    private readonly ConcurrentCircularBuffer<string> concurrentCircularBuffer;

    private Thread? communicationThread;

    private bool disposedValue;
    private bool isAlive;

    private string? activeSerialPortName;

    public SerialConnection() => concurrentCircularBuffer = new ConcurrentCircularBuffer<string>(5_000);

    public bool Connect(string portName)
    {
        if (activeSerialPortName != null)
        {
            return false;
        }

        activeSerialPortName = portName;

        communicationThread = new Thread(Run);
        communicationThread.Start();

        return communicationThread.IsAlive;
    }

    public async Task<bool> Disconnect()
    {
        if (activeSerialPortName == null || !isAlive)
        {
            return true;
        }

        var disconnectTask = new Task<bool>(() =>
        {
            isAlive = false;
            communicationThread?.Join();
            return communicationThread?.IsAlive == false;
        });

        disconnectTask.Start();
        return await disconnectTask;
    }

    public void Dispose()
    {
        Dispose(disposing: true);
        GC.SuppressFinalize(this);
    }

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

    private void Run()
    {
        if (activeSerialPortName == null)
        {
            return;
        }

        var port = new SerialPort(activeSerialPortName)
        {
            ReadTimeout = 500,
            WriteTimeout = 500
        };

        port.Open();

        isAlive = true;

        ulong pollCount = 0;

        while (isAlive)
        {
            try
            {
                concurrentCircularBuffer.Write(port.ReadLine());
                pollCount = 0;
            }
            catch (TimeoutException)
            {
                if (++pollCount >= KILL_CONNECTION_AFTER_N_EMPTY_POLLS)
                {
                    isAlive = false;
                }
            }
        }

        port.Close();
        activeSerialPortName = null;
    }
}
