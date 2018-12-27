package com.emc.ssmp.mybatis;

import org.junit.Test;

public class FinallyTest {

	@Test
	public void testFinalyMethod(){
		try{
			String[] str = new String[3];
			String split = str[1];
//			String split = str[5];
			System.out.println(1);
		}catch(Exception e){
			System.out.println(2);
		}finally{
			System.out.println(3);
		}
	}
	
}
