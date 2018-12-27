package com.emc.ssmp.service;

import org.apache.ibatis.session.SqlSessionFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import com.emc.ssmp.dao.UserDao;
import com.emc.ssmp.pojo.User;

@Service
public class UserService {

	@Autowired
	private UserDao userDao;
	
//	@Autowired
	public SqlSessionFactory  sqlSessionFactory;
	
/*  private UserMapper userMapper;

  public void setUserMapper(UserMapper userMapper) {
    this.userMapper = userMapper;
  }
*/
	
	public User getUserById(String userId) {
	  return userDao.getUser(userId);
	}
	
	public User getOne(){
		User u = new User();
		u.setUserid("shanghai");
		u.setAddress("emc");
		return u;
	}
	
}
