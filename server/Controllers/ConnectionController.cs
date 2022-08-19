using Microsoft.AspNetCore.Mvc;
using server.Models.Connection;
using server.Services;

namespace server.Controllers;

[ApiController]
[Route("connection")]
public class ConnectionController : ControllerBase
{
    private readonly SerialConnectionDiscoveryService connectionService;

    public ConnectionController(SerialConnectionDiscoveryService connectionService)
    {
        this.connectionService = connectionService;
    }

    [HttpGet("ports")]
    public ActionResult<AvailablePorts> ListAvailableSerialPorts()
    {
        return new AvailablePorts
        {
            PortNames = connectionService.GetAllAvailablePorts()
        };
    }

    [HttpPost("ports")]
    public EmptyResult SelectActivePort(SelectedPort selectedPort)
    {
        if (!connectionService.ChangePort(selectedPort.PortName))
        {
            Response.StatusCode = 406;
        }

        return new EmptyResult();
    }
}
