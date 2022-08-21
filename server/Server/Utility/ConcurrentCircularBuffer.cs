namespace Server.Utility;

public class ConcurrentCircularBuffer<T> where T : class
{
    private readonly ReaderWriterLockSlim bufferLock;
    private readonly T?[] buffer;

    private int readIndex;
    private int writeIndex;
    private int readWriteDelta;

    /// <summary>
    /// Create a new, thread-safe, <see href="https://en.wikipedia.org/wiki/Circular_buffer">circular buffer</see>.
    /// Please note that this implementation overrides old data if the reader is slower than the writer.
    /// </summary>
    /// <param name="capacity"></param>
    public ConcurrentCircularBuffer(int capacity)
    {
        bufferLock = new ReaderWriterLockSlim();
        buffer = new T?[capacity];

        readIndex = 0;
        writeIndex = 0;
    }

    /// <summary>
    /// Maximum number of items that can be stored in the buffer
    /// </summary>
    public int Capacity => buffer.Length;

    /// <summary>
    /// Write a single value to the buffer. Overrides old values if the reader is too slow.
    /// </summary>
    /// <param name="value">Value to write</param>
    public void Write(T? value)
    {
        bufferLock.EnterWriteLock();

        buffer[writeIndex++] = value;
        writeIndex %= buffer.Length;

        readWriteDelta++;

        bufferLock.ExitWriteLock();
    }

    /// <summary>
    /// Write multiple values to the buffer. Overrides old values if the reader is too slow.
    /// </summary>
    /// <param name="values">Values to write</param>
    public void WriteAll(params T?[] values)
    {
        bufferLock.EnterWriteLock();

        foreach (var value in values)
        {
            buffer[writeIndex++] = value;
            writeIndex %= buffer.Length;

            readWriteDelta++;
        }

        bufferLock.ExitWriteLock();
    }

    /// <summary>
    /// Read the next available data
    /// </summary>
    /// <returns>Data or null if no new data was available</returns>
    public T? Read()
    {
        T? value = null;
        bufferLock.EnterReadLock();

        // Only read if data is available
        if (readWriteDelta > 0)
        {
            value = buffer[readIndex++];
            readIndex %= buffer.Length;

            readWriteDelta--;
            readWriteDelta %= buffer.Length;
        }

        bufferLock.ExitReadLock();
        return value;
    }

    /// <summary>
    /// Read all available data in the buffer up until the writer
    /// </summary>
    /// <returns>Values read</returns>
    public List<T?> ReadAll()
    {
        var data = new List<T?>();
        bufferLock.EnterReadLock();

        while (readWriteDelta > 0)
        {
            data.Add(buffer[readIndex++]);
            readIndex %= buffer.Length;

            readWriteDelta--;
            readWriteDelta %= buffer.Length;
        }

        bufferLock.ExitReadLock();
        return data;
    }

    /// <summary>
    /// Gracefully disposes of the read / write lock
    /// </summary>
    ~ConcurrentCircularBuffer()
    {
        bufferLock?.Dispose();
    }
}
