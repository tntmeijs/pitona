namespace Server.Utility;

public class Obd2Connection : SerialConnection
{
    private const int RESPONSE_TIMEOUT_MS = 1_000;

    public Obd2Connection(ILogger<SerialConnection> logger) : base(logger)
    {
    }

    public void SendSingleObd2Command(string mode, string pid, string? data)
    {
        if (data == null)
        {
            SendRaw($"{mode}{pid}\r");
        }
        else
        {
            SendRaw($"{mode}{pid}{data}\r");
        }
    }

    public async Task<string?> WaitForResponse(int expectedBytes)
    {
        return await Task
            .Run(() => PollIndefinitelyForResponse(expectedBytes))
            .WaitAsync(TimeSpan.FromMilliseconds(RESPONSE_TIMEOUT_MS));
    }

    private string? PollIndefinitelyForResponse(int expectedBytes)
    {
        string? response = "";

        while (response.Length < expectedBytes)
        {
            var data = concurrentCircularBuffer.Read();

            if (data != null)
            {
                response += data;
            }
        }

        return response;
    }
}
