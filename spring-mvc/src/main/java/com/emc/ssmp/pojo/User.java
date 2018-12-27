package com.emc.ssmp.pojo;

import java.io.Serializable;

public class User implements Serializable{

	
	/**
	 * 
	 */
	private static final long serialVersionUID = 8380979239087169486L;
	private String userid;
	private String password;
	private String email;
	private String telephone;
	private String address;
	public String getUserid() {
		return userid;
	}
	public void setUserid(String userid) {
		this.userid = userid;
	}
	public String getPassword() {
		return password;
	}
	
	public void setPassword(String password) {
		this.password = password;
	}
	public String getEmail() {
		return email;
	}
	public void setEmail(String email) {
		this.email = email;
	}

	public String getTelephone() {
		return telephone;
	}
	public void setTelephone(String telephone) {
		this.telephone = telephone;
	}
	public String getAddress() {
		return address;
	}
	public void setAddress(String address) {
		this.address = address;
	}
	

	public User() {
		super();
	}
	public User(String userid, String password, String email, String telephone,
			String address) {
		super();
		this.userid = userid;
		this.password = password;
		this.email = email;
		this.telephone = telephone;
		this.address = address;
	}
	
}
