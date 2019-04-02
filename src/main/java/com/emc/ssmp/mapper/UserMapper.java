package com.emc.ssmp.mapper;

import com.emc.ssmp.pojo.User;

//@Repository
public interface UserMapper {

	
//	@Select("SELECT * FROM users WHERE userid = #{userId}")
//	User getUser(@Param("userId") String userId);
	
	public User getUser(String userId);
}
