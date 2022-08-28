namespace Server.Models.Obd2;

/// <summary>
/// OBD2 failure without any additional information
/// </summary>
public class Obd2EmptyFailureResult : IObd2Result
{
    public bool Success => false;
}

/// <summary>
/// OBD2 raw result
/// </summary>
public class Obd2RawDataResult : IObd2Result
{
    public bool Success => true;
    public string Data { get; set; } = string.Empty;
}
