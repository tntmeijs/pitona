package dev.tahar.pitona.system

import dev.tahar.spec.api.SystemApi
import dev.tahar.spec.model.SystemInformationResponseV1
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController
import java.time.Instant

@RestController
internal class SystemController(private val systemStatusService: SystemStatusService) : SystemApi {

    /**
     * Return system status information
     */
    override fun getSystemStatus(): ResponseEntity<SystemInformationResponseV1> {
        val response = SystemInformationResponseV1(
            systemStatusService.getCpuTemperature(),
            Instant.now().toEpochMilli()
        );

        return ResponseEntity.ok(response);
    }

}
