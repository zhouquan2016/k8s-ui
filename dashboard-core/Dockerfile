FROM maven:3.8.6-openjdk-18 as build
WORKDIR /app
COPY . .
RUN mvn clean compile package -Dmaven.test.skip=true -s ./maven-settings.xml

FROM openjdk:18-jdk-oraclelinux8 as deploy
WORKDIR /app
COPY --from=build /app/target/dashboard-core-1.0.0.jar /app
CMD ["java", "-jar", "/app/dashboard-core-1.0.0.jar"]