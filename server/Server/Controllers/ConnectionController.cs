using Microsoft.AspNetCore.Mvc;
using Server.Models.Connection;
using Server.Services;

namespace server.Controllers;

[ApiController]
[Route("connection")]
public class ConnectionController : ControllerBase
{
    private readonly Obd2ConnectionService serialConnectionService;

    public ConnectionController(Obd2ConnectionService serialConnectionService) => this.serialConnectionService = serialConnectionService;

    [HttpGet("ports")]
    public async Task<ActionResult<AvailablePorts>> ListAvailableSerialPorts()
    {
        return new AvailablePorts
        {
            PortNames = await serialConnectionService.GetAllAvailablePorts()
        };
    }

    [HttpPost("ports")]
    public async Task<ActionResult> SelectActivePort(SelectedPort selectedPort)
    {
        if (!await serialConnectionService.ChangePort(selectedPort.PortName))
        {
            return BadRequest();
        }

        return Ok();
    }

    [HttpDelete("ports")]
    public async Task<ActionResult> Disconnect()
    {
        var success = await serialConnectionService.Disconnect();

        if (!success)
        {
            return BadRequest();
        }

        return Ok();
    }
}
