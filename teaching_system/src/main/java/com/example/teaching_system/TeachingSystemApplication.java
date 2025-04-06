package com.example.teaching_evaluation;

import org.mybatis.spring.annotation.MapperScan;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.context.annotation.ComponentScan;

@SpringBootApplication
@MapperScan(value = "com.example.mapper.*")
@ComponentScan(value = "com.example.*")
public class TeachingSystemApplication {

    public static void main(String[] args) {
        SpringApplication.run(TeachingSystemApplication.class, args);
    }

}
