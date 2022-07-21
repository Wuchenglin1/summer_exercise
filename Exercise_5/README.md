# lv0

## Nginx

![image-20220720200009339](http://110.42.184.72:8092/1658318410.png)

## Apache

### 配置文件

```
ServerRoot "/etc/httpd"			#用于指定Apache运行的根目录
Listen 8070						#监听80端口
MaxClients  256					#指定同时能访问服务器的客户机数量为256
DocumentRoot "/root/index"		#网页文件存放的目录
DirectoryIndex index.html		#主页
Include conf.d/*.conf			#加载/etc/httpd/conf/conf.d/目录中所有以.conf结尾的文件
ServerName 110.42.184.72		#域名
Include conf.d/*.conf			#包含的子配置文件
User root						#用户是root
Group root						#用户组是root
```

### 效果

<img src="http://110.42.184.72:8092/1658325964.png" alt="image-20220720220603555" style="zoom:50%;" />

# Lv0.5

由于8081,8082,8083端口被占了，我就使用的8087,8088,8089端口

![image-20220720204631959](http://110.42.184.72:8092/1658321192.png)



LV1

在小[小网站](https://killercoda.com/playgrounds/scenario/kubernetes)上开k8s测试速度太感人了，到后面吗直接拉不下来了:sob: 

![image-20220721145342943](http://110.42.184.72:8092/1658386425.png)

~~然后换了个nginx开~~

![image-20220721150802469](http://110.42.184.72:8092/1658387284.png)

