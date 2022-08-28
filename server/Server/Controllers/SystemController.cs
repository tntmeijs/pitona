using Microsoft.AspNetCore.Mvc;
using Server.Models.System;

namespace server.Controllers;

[ApiController]
[Route("system")]
public class SystemController : ControllerBase
{
    [HttpGet("status")]
    public ActionResult<StatusModel> GetStatus()
    {
        return Ok(new StatusModel { CpuTemperature = 69420.0d });
    }
}
