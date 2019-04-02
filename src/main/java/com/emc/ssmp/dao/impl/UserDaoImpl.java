package com.emc.ssmp.dao.impl;



import javax.sql.DataSource;

import org.apache.ibatis.session.SqlSessionFactory;
import org.mybatis.spring.support.SqlSessionDaoSupport;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.jdbc.core.JdbcTemplate;
import org.springframework.stereotype.Repository;

import com.emc.ssmp.dao.UserDao;
import com.emc.ssmp.mapper.UserMapper;
import com.emc.ssmp.pojo.User;

@Repository
public class UserDaoImpl  extends SqlSessionDaoSupport  implements UserDao{

	//extends SqlSessionDaoSupport 
//	@Autowired(required = true)
//	public SqlSessionFactory  sqlSessionFactory;
	
	@Autowired
	  @Override
	  public void setSqlSessionFactory(SqlSessionFactory sqlSessionFactory) {
	    super.setSqlSessionFactory(sqlSessionFactory);
	  }
	
	@Autowired(required = true)
	public DataSource dataSource;
	
/*   @Autowired
   public void setDataSource(DataSource dataSource) {
       this.jdbcTemplate = new JdbcTemplate(dataSource);
      // this.dataSource = dataSource;
   }*/
	
//	@Autowired
	public UserMapper userMapper;
	
//	@Autowired
//	@Qualifier("jdbcTemplate")
	private JdbcTemplate jdbcTemplate;


	
	
	public User getUser(String userId) {
		// TODO Auto-generated method stub
		UserMapper mapper =  
	            getSqlSession().getMapper(UserMapper.class);  
		return mapper.getUser(userId);
		
//		jdbcTemplate.execute("select * from user");
		
//		try {
//			dataSource.getConnection();
//		} catch (SQLException e) {
//			// TODO Auto-generated catch block
//			e.printStackTrace();
//		}
//		
//		
//		return new User();
//		return null;
	}

}
