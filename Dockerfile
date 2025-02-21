FROM goviewfile-base:v1

# 设置固定的项目路径
ENV WORKDIR /var/www/GoViewFile

# 添加应用可执行文件，并设置执行权限
ADD go_view_file   $WORKDIR/
COPY fonts/* /usr/share/fonts/ChineseFonts/

RUN  cd $WORKDIR && chmod +x $WORKDIR/go_view_file && mkdir cache && cd cache && \
     mkdir convert download local pdf && yum clean all


# 添加I18N多语言文件、静态文件、配置文件、模板文件
ADD public/  $WORKDIR/public
ADD config/  $WORKDIR/config
ADD template/ $WORKDIR/template
# jar包，用于将.msg文件转eml文件
COPY library/emailconverter-2.5.3-all.jar   /usr/local/emailconverter-2.5.3-all.jar
#pdf 添加水印
COPY library/pdfcpu    /usr/local/bin/pdfcpu

WORKDIR $WORKDIR
# 如果需要进入容器调式，可以注释掉下面的CMD. 
CMD  ./go_view_file


# ------------------------------------本地打包镜像---------------------
# docker build -t  goviewfile:v0.7  .
# docker run -d  -p 8082:8082 镜像ID
