using Microsoft.AspNetCore.Mvc;
using server.Models.Connection;
using server.Services;

namespace server.Controllers;

[ApiController]
[Route("connection")]
public class ConnectionController : ControllerBase
{
    private readonly SerialConnectionDiscoveryService serialConnectionService;

    public ConnectionController(SerialConnectionDiscoveryService serialConnectionService) => this.serialConnectionService = serialConnectionService;

    [HttpGet("ports")]
    public ActionResult<AvailablePorts> ListAvailableSerialPorts()
    {
        return new AvailablePorts
        {
            PortNames = serialConnectionService.GetAllAvailablePorts()
        };
    }

    [HttpPost("ports")]
    public async Task<EmptyResult> SelectActivePort(SelectedPort selectedPort)
    {
        if (!await serialConnectionService.ChangePort(selectedPort.PortName))
        {
            Response.StatusCode = 409;
        }

        return new EmptyResult();
    }

    [HttpDelete("ports")]
    public async Task<EmptyResult> Disconnect()
    {
        var success = await serialConnectionService.Disconnect();

        if (!success)
        {
            Response.StatusCode = 409;
        }

        return new EmptyResult();
    }
}
