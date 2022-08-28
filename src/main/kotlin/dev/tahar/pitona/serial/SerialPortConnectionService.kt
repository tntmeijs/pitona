package dev.tahar.pitona.serial

import com.fazecast.jSerialComm.SerialPort
import org.springframework.stereotype.Service

@Service
internal class SerialPortConnectionService {

    /**
     * Return the names of all available serial ports
     */
    fun getAllAvailableSerialPortNames(): List<String> {
        val foundDevices = mutableListOf<String>();

        for (port in SerialPort.getCommPorts()) {
            foundDevices.add(port.systemPortName);
        }

        return foundDevices.sorted();
    }

}
