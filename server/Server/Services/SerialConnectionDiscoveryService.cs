using System.IO.Ports;
using Server.Utility;

namespace Server.Services;

public class SerialConnectionDiscoveryService
{
    private readonly ILogger<SerialConnectionDiscoveryService> logger;
    private readonly SerialConnection serialConnection;

    public SerialConnectionDiscoveryService(ILogger<SerialConnectionDiscoveryService> logger)
    {
        this.logger = logger;
        serialConnection = new SerialConnection();
    }

    public List<string> GetAllAvailablePorts()
    {
        var ports = SerialPort.GetPortNames().ToList();
        ports.Sort();

        return ports;
    }

    public async Task<bool> ChangePort(string portName)
    {
        if (!SerialPort.GetPortNames().Any(x => x == portName))
        {
            logger.LogError("Cannot change serial port because no serial port with this name exists");
            return false;
        }

        if (await serialConnection.Disconnect())
        {
            return serialConnection.Connect(portName);
        }

        logger.LogError("Unable to disconnect old serial connection");
        return false;
    }

    public Task<bool> Disconnect()
    {
        return serialConnection.Disconnect();
    }
}
