package dev.tahar.pitona.system

import org.springframework.stereotype.Service

/**
 * Responsible for retrieving vital system information
 */
@Service
internal class SystemStatusService {

    /**
     * Returns the CPU temperature in degrees Celsius
     * @return Temperature in degrees Celsius
     */
    internal fun getCpuTemperature(): Double {
        return Double.NaN;
    }

}
