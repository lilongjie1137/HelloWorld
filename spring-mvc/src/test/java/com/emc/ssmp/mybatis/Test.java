package com.emc.ssmp.mybatis;

import java.util.HashMap;
import java.util.Map;

/** 
* @author 作者: Roger.Li
* @version 创建时间：2017年6月28日 下午6:05:11 
* 类说明 :
*/
public class Test {

	public static void main(String[] args) {
//		String s = "config_file:C:\\software\\redis\\redis.windows.conf";
//		String key=s.substring(0,s.indexOf(":"));
//		String value=s.substring(s.indexOf(":")+1);
//		System.out.println(key+" "+value);
		System.out.println("serverRequests".toUpperCase());
		String sb = "Active connections: 211 server accepts handled requests 25 25 7132 Reading: 0 Writing: 172 Waiting: 0 ";
		String activeConnections =sb.substring(sb.indexOf(":")+2,sb.indexOf("server")-1);
		String serverRequests = sb.substring(sb.indexOf("requests")+9, sb.indexOf("Reading")-1).split(" ")[2];
		String Writing = sb.substring(sb.indexOf("Writing:")+9, sb.indexOf("Waiting")-1);
		System.out.println(activeConnections);
		System.out.println(serverRequests);
		System.out.println(Writing);
	}

}
 