package com.emc.ssmp.service;

import java.awt.Graphics2D;
import java.awt.Image;
import java.awt.geom.AffineTransform;
import java.awt.image.AffineTransformOp;
import java.awt.image.BufferedImage;
import java.io.File;
import java.util.ArrayList;
import java.util.Arrays;
import java.util.List;

import javax.imageio.ImageIO;

/** 
* @author 作者: Roger.Li
* @version 创建时间：2018年9月1日 上午11:06:06 
* 类说明 :
*/
public class ImageTest {
	
	public static void main(String[] args) throws Exception {
//		BufferedImage finallImg1 = new ImageTest().mergeDown("C:\\Users\\Administrator\\Desktop\\base\\", 1, 10);
//		ImageIO.write(finallImg1,"jpg",new File("C:\\Users\\Administrator\\Desktop\\base\\world1.jpg"));
//		BufferedImage finallImg2 = new ImageTest().mergeDown("C:\\Users\\Administrator\\Desktop\\base\\", 11, 20);
//		ImageIO.write(finallImg2,"jpg",new File("C:\\Users\\Administrator\\Desktop\\base\\world2.jpg"));
//		BufferedImage finallImg3 = new ImageTest().mergeDown("C:\\Users\\Administrator\\Desktop\\base\\", 21, 30);
//		ImageIO.write(finallImg3,"jpg",new File("C:\\Users\\Administrator\\Desktop\\base\\world3.jpg"));
//		BufferedImage finallImg4 = new ImageTest().mergeDown("C:\\Users\\Administrator\\Desktop\\base\\", 31, 40);
//		ImageIO.write(finallImg4,"jpg",new File("C:\\Users\\Administrator\\Desktop\\base\\world4.jpg"));
		
		BufferedImage finallImg5 = new ImageTest().mergeDown("C:\\Users\\Administrator\\Desktop\\base\\", 41, 50);
		ImageIO.write(finallImg5,"jpg",new File("C:\\Users\\Administrator\\Desktop\\base\\world5.jpg"));

		
		//		BufferedImage bufferedImage = ImageIO.read(new File("C:\\Users\\Administrator\\Desktop\\40.jpg"));
//		BufferedImage fi = new ImageTest().zoomImage(bufferedImage, 580, 1031);
//		ImageIO.write(fi,"jpg",new File("C:\\Users\\Administrator\\Desktop\\40-2.jpg"));
	}

	public BufferedImage mergeRight(File file) throws Exception {
		int w = 160;
		int h = 205;
		
		BufferedImage bufferedImage = ImageIO.read(file);
		BufferedImage fa = bufferedImage.getSubimage(180,108,218,222);
		BufferedImage s1 = bufferedImage.getSubimage(210,388,w,h);
		BufferedImage s2 = bufferedImage.getSubimage(395,388,w,h);
		BufferedImage s3 = bufferedImage.getSubimage(23,620,w,h);
		BufferedImage s4 = bufferedImage.getSubimage(210,620,w,h);
		BufferedImage s5 = bufferedImage.getSubimage(395,620,w,h);
		BufferedImage faz = zoomImage(fa, w, h);
		
		List<BufferedImage> images = new ArrayList<>();
		images.add(faz);
		images.add(s1);
		images.add(s2);
		images.add(s3);
		images.add(s4);
		images.add(s5);
		List<BufferedImage> images_rank = rankStar(images);
		 
		BufferedImage finallImg = mergeImage(images_rank, 1);
        
		return finallImg;
	}
	
	public static List<BufferedImage> rankStar(List<BufferedImage> images_random) throws Exception {
		BufferedImage[] images = new BufferedImage[images_random.size()];
		for (int i = 0; i < images_random.size(); i++) {
			BufferedImage bi = images_random.get(i);
			if (i==0) {
				images[0] = bi;
			}else {
				images[judgeStar(bi)] = bi;
			}
		}
		return Arrays.asList(images);
	}
	
	public static int judgeStar(BufferedImage s) throws Exception {
		int star = 0;
		int[] rgbArray = new int[25];
		int[] rgb = s.getRGB(45, 115, 5, 5, rgbArray,0, 5);
		float hit = 0;
		for (int i : rgb) {
			int r = (i & 16711680) >> 16;
			int g = (i & 65280) >> 8;
			int b = (i & 255);
			if (r>200) {
				hit++;
			};
		}
		if (hit/25>0.8) {
			star = 5;
		}else {
			rgbArray = new int[25];
			rgb = s.getRGB(52, 115, 5, 5, rgbArray,0, 5);
			hit = 0;
			for (int i : rgb) {
				int r = (i & 16711680) >> 16;
				if (r>200) {
					hit++;
				};
			}
			if (hit/25>0.8) {
				star = 4;
			}else {
				rgbArray = new int[25];
				rgb = s.getRGB(60, 115, 5, 5, rgbArray,0, 5);
				hit = 0;
				for (int i : rgb) {
					int r = (i & 16711680) >> 16;
					if (r>200) {
						hit++;
					};
				}
				if (hit/25>0.8) {
					star = 3;
				}else {
					rgbArray = new int[25];
					rgb = s.getRGB(70, 115, 5, 5, rgbArray,0, 5);
					hit = 0;
					for (int i : rgb) {
						int r = (i & 16711680) >> 16;
						if (r>200) {
							hit++;
						};
					}
					if (hit/25>0.8) {
						star = 2;
					}else {
						rgbArray = new int[25];
						rgb = s.getRGB(78, 115, 5, 5, rgbArray,0, 5);
						hit = 0;
						for (int i : rgb) {
							int r = (i & 16711680) >> 16;
							if (r>200) {
								hit++;
							};
						}
						if (hit/25>0.8) {
							star = 1;
						}else {
							throw new Exception("no match star");
						}
					}
				}
			}
		}
		return star;
	}
	
	
	
	public BufferedImage mergeDown(String path,int str, int end) throws Exception {
		List<BufferedImage> images = new ArrayList<>();
		for (int i = str; i <= end; i++) {
			images.add(mergeRight(new File(path+i+".jpg")));
		}
		BufferedImage finallImg = mergeImage(images, 2);
		return finallImg;
	}
	
    /**
     * 图片缩放
     * @param fa 图片源文件
     * @param w 目标宽度
     * @param h 目标高度
     * @return 
     * @throws Exception
     */
    public static BufferedImage zoomImage(BufferedImage fa, int w,int h) throws Exception {
    	Image Itemp = fa.getScaledInstance(w, h, fa.SCALE_SMOOTH);//设置缩放目标图片模板
		double wr=w*1.0/fa.getWidth();     //获取缩放比例
		double hr=h*1.0/fa.getHeight();

        AffineTransformOp ato = new AffineTransformOp(AffineTransform.getScaleInstance(wr, hr), null);
        Itemp = ato.filter(fa, null);
        return (BufferedImage) Itemp;
    }
	
 
	/**
	 * @Description:图片拼接 （注意：必须两张图片长宽一致哦）
	 * @author:liuyc
	 * @time:2016年5月27日 下午5:52:24
	 * @param files 要拼接的文件列表
	 * @param type:1  横向拼接， 2 纵向拼接
	 */
	public static BufferedImage mergeImage(List<BufferedImage> images, int type) {
		int len = images.size();
		if (len < 1) {
			throw new RuntimeException("图片数量小于1");
		}
		int[][] ImageArrays = new int[len][];
		for (int i = 0; i < len; i++) {
			int width = images.get(i).getWidth();
			int height = images.get(i).getHeight();
			ImageArrays[i] = new int[width * height];
			ImageArrays[i] = images.get(i).getRGB(0, 0, width, height, ImageArrays[i], 0, width);
		}
		int newHeight = 0;
		int newWidth = 0;
		for (int i = 0; i < images.size(); i++) {
			// 横向
			if (type == 1) {
				newHeight = newHeight > images.get(i).getHeight() ? newHeight : images.get(i).getHeight();
				newWidth += images.get(i).getWidth();
			} else if (type == 2) {// 纵向
				newWidth = newWidth > images.get(i).getWidth() ? newWidth : images.get(i).getWidth();
				newHeight += images.get(i).getHeight();
			}
		}
//		if (type == 1 && newWidth < 1) {
//			return;
//		}
//		if (type == 2 && newHeight < 1) {
//			return;
//		}
 
		// 生成新图片
		try {
			BufferedImage ImageNew = new BufferedImage(newWidth, newHeight, BufferedImage.TYPE_INT_RGB);
			int height_i = 0;
			int width_i = 0;
			for (int i = 0; i < images.size(); i++) {
				if (type == 1) {
					ImageNew.setRGB(width_i, 0, images.get(i).getWidth(), newHeight, ImageArrays[i], 0,
							images.get(i).getWidth());
					width_i += images.get(i).getWidth();
				} else if (type == 2) {
					ImageNew.setRGB(0, height_i, newWidth, images.get(i).getHeight(), ImageArrays[i], 0, newWidth);
					height_i += images.get(i).getHeight();
				}
			}
			//输出想要的图片
//			ImageIO.write(ImageNew, targetFile.split("\\.")[1], new File(targetFile));
			return ImageNew;
 
		} catch (Exception e) {
			throw new RuntimeException(e);
		}
	}
 
	/**
	 * @Description:小图片贴到大图片形成一张图(合成)
	 * @author:liuyc
	 * @time:2016年5月27日 下午5:51:20
	 */
	public static final void overlapImage(String bigPath, String smallPath, String outFile) {
		try {
			BufferedImage big = ImageIO.read(new File(bigPath));
			BufferedImage small = ImageIO.read(new File(smallPath));
			Graphics2D g = big.createGraphics();
			int x = (big.getWidth() - small.getWidth()) / 2;
			int y = (big.getHeight() - small.getHeight()) / 2;
			g.drawImage(small, x, y, small.getWidth(), small.getHeight(), null);
			g.dispose();
			ImageIO.write(big, outFile.split("\\.")[1], new File(outFile));
		} catch (Exception e) {
			throw new RuntimeException(e);
		}
	}
	
	/** 
	    * Java拼接多张图片 
	    *  
	    * @param imgs 
	    * @param type 
	    * @param dst_pic 
	    * @return 
	    */  
	  public static boolean merge(String[] imgs, String type, String dst_pic) {  
		    //获取需要拼接的图片长度
		       int len = imgs.length;  
		       //判断长度是否大于0
		       if (len < 1) {  
		           return false;  
		       }  
		       File[] src = new File[len];  
		       BufferedImage[] images = new BufferedImage[len];  
		       int[][] ImageArrays = new int[len][];  
		       for (int i = 0; i < len; i++) {  
		           try {  
		               src[i] = new File(imgs[i]);  
		               images[i] = ImageIO.read(src[i]);  
		           } catch (Exception e) {  
		               e.printStackTrace();  
		               return false;  
		           }  
		           int width = images[i].getWidth();  
		           int height = images[i].getHeight();  
		        // 从图片中读取RGB 像素
		           ImageArrays[i] = new int[width * height];
		           ImageArrays[i] = images[i].getRGB(0, 0, width, height,  ImageArrays[i], 0, width);  
		       }  
		 
		       int dst_height = 0;  
		       int dst_width = images[0].getWidth();  
		     //合成图片像素
		       for (int i = 0; i < images.length; i++) {  
		           dst_width = dst_width > images[i].getWidth() ? dst_width     : images[i].getWidth();  
		           dst_height += images[i].getHeight();  
		       }  
		       //合成后的图片
		       System.out.println("宽度:"+dst_width);  
		       System.out.println("高度:"+dst_height);  
		       if (dst_height < 1) {  
		           System.out.println("dst_height < 1");  
		           return false;  
		       }  
		       // 生成新图片   
		       try {  
		           // dst_width = images[0].getWidth();   
		           BufferedImage ImageNew = new BufferedImage(dst_width, dst_height,  
		                   BufferedImage.TYPE_INT_RGB);  
		           int height_i = 0;  
		           for (int i = 0; i < images.length; i++) {  
		               ImageNew.setRGB(0, height_i, dst_width, images[i].getHeight(),  
		                       ImageArrays[i], 0, dst_width);  
		               height_i += images[i].getHeight();  
		           }  
		 
		           File outFile = new File(dst_pic);  
		           ImageIO.write(ImageNew, type, outFile);// 写图片 ，输出到硬盘 
		       } catch (Exception e) {  
		           e.printStackTrace();  
		           return false;  
		       }  
		       return true;  
		   }  
		
}
 