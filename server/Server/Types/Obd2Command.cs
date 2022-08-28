namespace Server.Types;

public enum Obd2Command
{
    /// =========================== ///
    /// Mode 01 - Show current data ///
    /// =========================== ///

    TIME_RUN_WITH_CHECK_ENGINE_LIGHT_ON,
    DISTANCE_TRAVELED_WITH_CHECK_ENGINE_LIGHT_ON
}

public static class Obd2CommandExtensions
{
    public static string Mode(this Obd2Command command)
    {
        return command switch
        {
            Obd2Command.TIME_RUN_WITH_CHECK_ENGINE_LIGHT_ON or
            Obd2Command.DISTANCE_TRAVELED_WITH_CHECK_ENGINE_LIGHT_ON => "01",
            _ => throw new NotImplementedException(),
        };
    }

    public static string Pid(this Obd2Command command)
    {
        return command switch
        {
            Obd2Command.TIME_RUN_WITH_CHECK_ENGINE_LIGHT_ON => "4D",
            Obd2Command.DISTANCE_TRAVELED_WITH_CHECK_ENGINE_LIGHT_ON => "21",
            _ => throw new NotImplementedException(),
        };
    }

    public static int ExpectedBytes(this Obd2Command command)
    {
        return command switch
        {
            Obd2Command.TIME_RUN_WITH_CHECK_ENGINE_LIGHT_ON or
            Obd2Command.DISTANCE_TRAVELED_WITH_CHECK_ENGINE_LIGHT_ON => 2,
            _ => throw new NotImplementedException()
        };
    }
}
