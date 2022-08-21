using Server.Utility;

namespace ServerTests;

public class ConcurrentCircularBufferTests
{
    private ConcurrentCircularBuffer<string> circularBuffer;

    public ConcurrentCircularBufferTests()
    {
        circularBuffer = new ConcurrentCircularBuffer<string>(5);
    }

    [Fact]
    public void ConstructorExpectInitialized()
    {
        Assert.Equal(5, circularBuffer.Capacity);
        Assert.Null(circularBuffer.Read());
        Assert.Empty(circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteReadExpectCorrectValue()
    {
        circularBuffer.Write("abc");
        Assert.Equal("abc", circularBuffer.Read());
    }

    [Fact]
    public void WriteReadExpectCorrectValues()
    {
        circularBuffer.Write("1");
        circularBuffer.Write("2");
        circularBuffer.Write("3");
        circularBuffer.Write("4");
        circularBuffer.Write("5");

        Assert.Equal("1", circularBuffer.Read());
        Assert.Equal("2", circularBuffer.Read());
        Assert.Equal("3", circularBuffer.Read());
        Assert.Equal("4", circularBuffer.Read());
        Assert.Equal("5", circularBuffer.Read());
    }

    [Fact]
    public void WriteAllReadExpectCorrectValues()
    {
        circularBuffer.WriteAll("1", "2", "3", "4", "5");

        Assert.Equal("1", circularBuffer.Read());
        Assert.Equal("2", circularBuffer.Read());
        Assert.Equal("3", circularBuffer.Read());
        Assert.Equal("4", circularBuffer.Read());
        Assert.Equal("5", circularBuffer.Read());
    }

    [Fact]
    public void WriteReadAllExpectCorrectValues()
    {
        circularBuffer.Write("1");
        circularBuffer.Write("2");
        circularBuffer.Write("3");
        circularBuffer.Write("4");
        circularBuffer.Write("5");

        Assert.Equal(new List<string?>() { "1", "2", "3", "4", "5" }, circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteAllReadAllExpectCorrectValues()
    {
        circularBuffer.WriteAll("1", "2", "3", "4", "5");
        Assert.Equal(new List<string?>() { "1", "2", "3", "4", "5" }, circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteAllPlusOneReadExpectReadOneValue()
    {
        // Write capacity + 1
        circularBuffer.WriteAll("1", "2", "3", "4", "5", "6");

        // Since the old data has been invalidated, only expect the reader to read one item
        Assert.Equal("6", circularBuffer.Read());
        Assert.Null(circularBuffer.Read());
    }

    [Fact]
    public void WriteAllPlusOneReadAllExpectReadOneValue()
    {
        // Write capacity + 1
        circularBuffer.WriteAll("1", "2", "3", "4", "5", "6");

        // Since the old data has been invalidated, only expect the reader to read one item
        Assert.Equal(new List<string?>() { "6" }, circularBuffer.ReadAll());
        Assert.Equal(new List<string?>(), circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteAllWithWrapAroundExpectFirstItemTruncated()
    {
        circularBuffer.WriteAll("1", "2", "3", "4", "5", "6", "7", "8", "9");
        Assert.Equal(new List<string?>() { "6", "7", "8", "9" }, circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteAllWithFullWrapAroundExpectFirstItemsTruncated()
    {
        circularBuffer.WriteAll("1", "2", "3", "4", "5", "6", "7", "8", "9", "10");
        Assert.Equal(new List<string?>() { "6", "7", "8", "9", "10" }, circularBuffer.ReadAll());
    }

    [Fact]
    public void WriteAllWithWrapAroundExpectSuccess()
    {
        circularBuffer.WriteAll("1", "2", "3", "4", "5", "6", "7", "8", "9", "10");
        Assert.Equal("6", circularBuffer.Read());

        circularBuffer.Write("11");
        Assert.Equal(new List<string?>() { "7", "8", "9", "10", "11" }, circularBuffer.ReadAll());
    }
}