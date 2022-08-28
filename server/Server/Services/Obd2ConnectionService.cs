using System.IO.Ports;
using Server.Models.Obd2;
using Server.Types;
using Server.Utility;

namespace Server.Services;

public class Obd2ConnectionService
{
    private readonly ILogger<Obd2ConnectionService> logger;
    private readonly Obd2Connection obd2Connection;

    public Obd2ConnectionService(ILogger<Obd2ConnectionService> logger, ILogger<SerialConnection> serialConnectionLogger)
    {
        this.logger = logger;
        obd2Connection = new Obd2Connection(serialConnectionLogger);
    }

    public async Task<List<string>> GetAllAvailablePorts()
    {
        return await Task.Run(() =>
        {
            var ports = SerialPort.GetPortNames().ToList();
            ports.Sort();

            return ports;
        });
    }

    public async Task<bool> ChangePort(string portName)
    {
        if (!SerialPort.GetPortNames().Any(x => x == portName))
        {
            logger.LogError("Cannot change serial port because no serial port with this name exists");
            return false;
        }

        return await Task.Run(() =>
        {
            if (obd2Connection.Disconnect())
            {
                return obd2Connection.Connect(portName);
            }

            logger.LogError("Unable to disconnect old serial connection");
            return false;
        });
    }

    public async Task<IObd2Result> Send(Obd2Command command, string? data)
    {
        obd2Connection.SendSingleObd2Command(command.Mode(), command.Pid(), data);
        //var response = await obd2Connection.WaitForResponse(command.ExpectedBytes());

        //if (response == null)
        //{
        //    return new Obd2EmptyFailureResult();
        //}

        //return new Obd2RawDataResult { Data = response };

        return new Obd2EmptyFailureResult();
    }

    public async Task<List<string?>> DumpData() => await Task.Run(() => obd2Connection.Dump());

    public async Task<bool> Disconnect() => await Task.Run(() => obd2Connection.Disconnect());
}
