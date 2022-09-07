package com.zhqn.dashboard.core.controller;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

/**
 * Copyright (C), 2015-2022, 易车
 * FileName: IndexController
 * Author:   zhouquan3
 * Date:     2022/9/7 16:14
 * Description:
 * @author zhouquan3
 */
@RestController
@RequestMapping("/")
public class IndexController {

    @GetMapping("/")
    public String index() {
        return "ok";
    }
}
