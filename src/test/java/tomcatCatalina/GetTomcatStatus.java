package tomcatCatalina;

import java.io.BufferedReader;
import java.io.IOException;
import java.io.InputStream;
import java.io.InputStreamReader;
import java.net.HttpURLConnection;
import java.net.MalformedURLException;
import java.net.URL;

import org.dom4j.Document;
import org.dom4j.DocumentException;
import org.dom4j.DocumentHelper;
import org.dom4j.Element;
import org.junit.Test;

/** 
* @author 作者: Roger.Li
* @version 创建时间：2017年7月4日 上午11:12:08 
* 类说明 :
*/
public class GetTomcatStatus {

	/**
	  * @Description: 获取指定URL的内容
	  * @Version1.0 2014-7-23 下午02:18:22 by xuqiang（xuqiang@merit.com）创建
	  * @param tempurl url地址
	  * @param username tomcat 管理用户名
	  * @param password tomcat 管理用户密码
	  * @return
	  * @throws IOException
	  */
	public static String getHtmlContext(String tempurl, String username, String password) throws IOException {
		URL url = null;
		BufferedReader breader = null;
		InputStream is = null;
		StringBuffer resultBuffer = new StringBuffer();
		try {
			url = new URL(tempurl);
			String userPassword = username + ":" + password;
			String encoding = new sun.misc.BASE64Encoder().encode (userPassword.getBytes());//在classpath中添加rt.jar包，在%java_home%/jre/lib/rt.jar

			HttpURLConnection conn = (HttpURLConnection) url.openConnection();
			conn.setRequestProperty("Authorization", "Basic " + encoding);
			is = conn.getInputStream();
			breader = new BufferedReader(new InputStreamReader(is));
			String line = "";
			while ((line = breader.readLine()) != null) {
				resultBuffer.append(line);
			}
		} catch (MalformedURLException e) {
			e.printStackTrace();
		} finally {
			if(breader != null)
				breader.close();
			if(is != null)
				is.close();
		}
		return resultBuffer.toString();
	}
	@Test
	public void test() {
		String result = "";
		Document doc = null;//引入org.dom4j包
		try {
			result = this.getHtmlContext("http://192.168.130.130:8080/manager/status?XML=true", "admin", "admin");
			doc = DocumentHelper.parseText(result);//将字符串转化为XML的Document
		} catch (IOException e) {
			e.printStackTrace();
		} catch (DocumentException e) {
			e.printStackTrace();
		}
		System.out.println(doc.asXML());
		
		if(doc != null ){
			Element root = doc.getRootElement();
			Element memory = (Element)root.selectSingleNode("/status/jvm/memory");
			String free = memory.attributeValue("free");
			System.out.println(free);
		}
	}
	
	
}
 