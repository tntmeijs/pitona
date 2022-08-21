namespace Server.Utility;

public class ConcurrentCircularBuffer<T> where T : class
{
    private readonly ReaderWriterLockSlim bufferLock;
    private readonly T?[] buffer;

    private int readIndex;
    private int writeIndex;

    public ConcurrentCircularBuffer(int capacity)
    {
        bufferLock = new ReaderWriterLockSlim();
        buffer = new T[capacity];

        readIndex = 0;
        writeIndex = 0;
    }

    public void Write(T? value)
    {
        bufferLock.EnterWriteLock();

        buffer[writeIndex++] = value;
        writeIndex %= buffer.Length;

        // If the reader is too slow, override old values
        if (readIndex < writeIndex)
        {
            readIndex = writeIndex;
        }

        bufferLock.ExitWriteLock();
    }

    public T? Read()
    {
        T? value = null;

        // Only read if data is available
        if (readIndex <= writeIndex)
        {
            bufferLock.EnterReadLock();

            value = buffer[readIndex++];

            bufferLock.ExitReadLock();
        }

        return value;
    }

    public List<T?> ReadAll()
    {
        var data = new List<T?>();
        T? value;

        while ((value = Read()) != null)
        {
            data.Add(value);
        }

        return data;
    }

    ~ConcurrentCircularBuffer()
    {
        bufferLock?.Dispose();
    }
}
