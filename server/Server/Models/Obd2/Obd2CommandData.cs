using Server.Types;

namespace Server.Models.Obd2;

public class Obd2CommandData
{
    public Obd2Command Type { get; set; }
    public string? Payload { get; set; } = null;
}
