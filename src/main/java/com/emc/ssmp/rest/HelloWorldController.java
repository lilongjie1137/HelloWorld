package com.emc.ssmp.rest;

import java.util.ArrayList;
import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Qualifier;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.ResponseBody;
import org.springframework.web.bind.annotation.RestController;

import com.emc.ssmp.pojo.User;
import com.emc.ssmp.service.UserService;
import com.google.gson.Gson;


@RestController
@RequestMapping(value="/helloWorld")
public class HelloWorldController {

	

	@Autowired
	private UserService userService;
	
	
	@RequestMapping(value="/one",method=RequestMethod.GET)
	@ResponseBody
	public String getAllSi(@RequestParam(value="name") String name){
		System.out.println("this is from controller:" + name);
		return "this is from controller perseus";
	}
	
	@RequestMapping(value="/user/{id}",method=RequestMethod.GET)
	@ResponseBody
	public String getOneUser(){
		User u = new User();
		u.setUserid("yanjun");
		u.setAddress("shanghai");
		return new Gson().toJson(u);
	}
	
	@RequestMapping(value="/user",method=RequestMethod.GET)
	@ResponseBody
	public String getUser(){
		List<User> uList = new ArrayList<User>();
		for(int i=0;i<5;i++){
			User u = new User(); 
			u.setUserid("perseus::" + i);
			u.setAddress("emc::" + i);
			uList.add(u);
		}
		return new Gson().toJson(uList);
	}
	
	@RequestMapping(value="/two",method=RequestMethod.GET)
	@ResponseBody
	public String getJsonStr(){
		return "{\"username\":\"perseus for json string\"}";
	}
	
	
	@RequestMapping(value="/userObject",method=RequestMethod.GET)
	public @ResponseBody User getUserObject(){
		User u = new User();
		u.setUserid("perseus");
		u.setAddress("emc");
		return u;
	}
	
	@RequestMapping(value="/mybatis/user",method=RequestMethod.GET)
	public @ResponseBody User getMybatisUser(@RequestParam(value="name") String name){
		User u = userService.getUserById(name);
//		User u = userService.getOne();
		return u;
	}
	
}
