#基于的基础镜像
FROM python:3.10.0
#代码添加到code文件夹
ADD . .
RUN chmod 111 doevaluate.sh
RUN chmod 000 Dockerfile
RUN chmod 111 evaluation.py
RUN chmod 444 test_dataset_5.pickle
WORKDIR .
# 安装支持
RUN pip3 config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple/
RUN pip3 install torch --no-cache-dir
RUN pip3 install numpy --no-cache-dir
RUN pip3 install keras --no-cache-dir
RUN pip3 install tensorflow --no-cache-dir
RUN pip3 install torchvision --no-cache-dir
CMD ["bash", "./doevaluate.sh"]
