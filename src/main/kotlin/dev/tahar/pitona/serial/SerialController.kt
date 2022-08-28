package dev.tahar.pitona.serial

import dev.tahar.spec.api.SerialApi
import dev.tahar.spec.model.GetAllSerialPortNames200Response
import org.springframework.http.ResponseEntity
import org.springframework.web.bind.annotation.RestController

@RestController
internal class SerialController(private val serialPortConnectionService: SerialPortConnectionService) : SerialApi {

    /**
     * Return the names of all available serial ports
     */
    override fun getAllSerialPortNames(): ResponseEntity<GetAllSerialPortNames200Response> {
        return ResponseEntity.ok(GetAllSerialPortNames200Response(serialPortConnectionService.getAllAvailableSerialPortNames()));
    }

}
