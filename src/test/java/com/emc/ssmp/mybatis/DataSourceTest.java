package com.emc.ssmp.mybatis;

import java.sql.Connection;
import java.sql.PreparedStatement;
import java.sql.ResultSet;
import java.sql.Statement;

import javax.sql.DataSource;

import org.springframework.context.ApplicationContext;
import org.springframework.context.support.ClassPathXmlApplicationContext;
import org.springframework.jdbc.core.JdbcTemplate;

public class DataSourceTest {

	ApplicationContext context = new ClassPathXmlApplicationContext(new String[] {"spring.xml"});
	
	public void testJdbcTemplate(){
		JdbcTemplate jt= (JdbcTemplate)context.getBean("jdbcTemplate");
		int count = jt.queryForInt("select count(*) from user");
		System.out.println("total count number is:" + count);
	}
	
	public void testStatement(){
		try{
			DataSource ds = (DataSource)context.getBean("dataSource");
			Connection con = ds.getConnection();
			Statement st = con.createStatement();
			PreparedStatement  ps = con.prepareStatement("select * from user");
			ResultSet rs = ps.executeQuery();
			while(rs.next()){
				System.out.println("rs result first row:"+rs.getString(1));
			}
		}catch(Exception e){
			e.printStackTrace();
		}
	}
	
	public static void main(String[] args)throws Exception{
		
		DataSourceTest t = new DataSourceTest();
//		t.testJdbcTemplate();
		t.testStatement();
		
	}
	
}
