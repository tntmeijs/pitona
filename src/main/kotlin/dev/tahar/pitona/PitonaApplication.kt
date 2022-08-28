package dev.tahar.pitona

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class PitonaApplication

fun main(args: Array<String>) {
	runApplication<PitonaApplication>(*args)
}
