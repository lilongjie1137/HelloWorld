package com.emc.ssmp.mybatis;

import java.util.logging.Level;
import java.util.logging.Logger;

import org.springframework.beans.BeansException;
import org.springframework.context.ApplicationContext;
import org.springframework.context.ApplicationContextAware;


public class SpringContextHelper implements ApplicationContextAware {
	private static final Logger logger = Logger.getLogger(SpringContextHelper.class.getName());

    private static ApplicationContext context ;
     
    /*
     * 注入ApplicationContext, 在加载Spring时自动获得context
     */
    public void setApplicationContext(ApplicationContext context)
            throws BeansException {
        SpringContextHelper.context = context;
        logger.log(Level.INFO, "set ApplicationContext=" + context.toString());
    }
     
    public static Object getBean(String beanName){
    	if (context == null)
    		return null;
        return context.getBean(beanName);
    }
}