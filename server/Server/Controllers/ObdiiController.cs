using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Server.Models.Obd2;
using Server.Services;
using Server.Types;

namespace Server.Controllers;

[ApiController]
[Route("obd2")]
public class Obd2Controller : ControllerBase
{
    private readonly ILogger<Obd2Controller> logger;
    private readonly Obd2ConnectionService serialConnectionService;

    public Obd2Controller(ILogger<Obd2Controller> logger, Obd2ConnectionService serialConnectionService)
    {
        this.logger = logger;
        this.serialConnectionService = serialConnectionService;
    }

    [HttpGet("dump")]
    public async Task<ActionResult<DataDump>> DumpBuffer()
    {
        var dump = new DataDump { RawData = await serialConnectionService.DumpData() };
        return Ok(dump);
    }

    [HttpPost("execute")]
    public async Task<ActionResult<IObd2Result>> ExecuteCommand(Obd2CommandData commandData)
    {
        var result = await serialConnectionService.Send(commandData.Type, commandData.Payload);

        if (result.Success)
        {
            return Ok(result);
        }

        return BadRequest();
    }
}
