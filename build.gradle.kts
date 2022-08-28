import org.jetbrains.kotlin.gradle.tasks.KotlinCompile

plugins {
    id("org.springframework.boot") version "2.7.3"
    id("io.spring.dependency-management") version "1.0.13.RELEASE"
    id("org.openapi.generator") version "6.0.1"
    kotlin("jvm") version "1.6.21"
    kotlin("plugin.spring") version "1.6.21"
}

group = "dev.tahar"
version = "0.0.1-SNAPSHOT"
java.sourceCompatibility = JavaVersion.VERSION_17

repositories {
    mavenCentral()
}

dependencies {
    implementation("org.springframework.boot:spring-boot-starter-web")
    implementation("com.fasterxml.jackson.module:jackson-module-kotlin")
    implementation("org.jetbrains.kotlin:kotlin-reflect")
    implementation("org.jetbrains.kotlin:kotlin-stdlib-jdk8")
    implementation("jakarta.validation:jakarta.validation-api:3.0.2")
    implementation("com.fazecast:jSerialComm:2.9.2")

    testImplementation("org.springframework.boot:spring-boot-starter-test")
}

openApiGenerate {
    generatorName.set("kotlin-spring")
    inputSpec.set("$rootDir/src/main/resources/static/api.yaml")
    apiPackage.set("dev.tahar.spec.api")
    modelPackage.set("dev.tahar.spec.model")
    outputDir.set("$buildDir/generated")

    configOptions.putAll(
        mapOf(
            "documentationProvider" to "none",
            "exceptionHandler" to "false",
            "interfaceOnly" to "true",
            "useBeanValidation" to "false",
            "useTags" to "true"
        )
    );
}

configure<SourceSetContainer> {
    named("main") {
        java.srcDir("$buildDir/generated/src/main/kotlin")
    }
}

tasks.named("compileKotlin") {
    dependsOn("openApiGenerate")
}

tasks.withType<KotlinCompile> {
    kotlinOptions {
        freeCompilerArgs = listOf("-Xjsr305=strict")
        jvmTarget = "17"
    }
}

tasks.withType<Test> {
    useJUnitPlatform()
}
