using System.IO.Ports;

namespace server.Services;

public class SerialConnectionDiscoveryService
{
    private readonly ILogger<SerialConnectionDiscoveryService> logger;

    public SerialConnectionDiscoveryService(ILogger<SerialConnectionDiscoveryService> logger)
    {
        this.logger = logger;
    }

    public List<string> GetAllAvailablePorts()
    {
        var ports = SerialPort.GetPortNames().ToList();
        ports.Sort();

        return ports;
    }

    public bool ChangePort(string portName)
    {
        if (!SerialPort.GetPortNames().Any(x => x == portName))
        {
            logger.LogError("Cannot change serial port because no serial port with this name exists");
            return false;
        }

        // TODO: implement
        logger.LogDebug("Attemping to change port...");
        return true;
    }
}
