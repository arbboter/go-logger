go-logger ��golang ����־�� �����ڶ�golang����log�ķ�װ��
�÷�����java��־���߰�log4j

��ӡ��־��6������ Debug��Info��Warn��Error��Fatal��Key��־�����ɵ͵��ߣ���Ӧ��֧��
Debugf��Infof��Warnf��Errorf��Fatalf��Keyf��ʽ���ķ�ʽ 

��־���ԣ�
1.�ο���https://github.com/donnie4w/go-logger
2.֧����־������������־��ĩ�ļ�λ��
3.��װ��һ�㣬��ֱ��ʹ��

ʾ����
    ��ʼ����logger.Init("logger", logger.DEBUG) // ָ����־ģ����+��־�ȼ�

��ӡ��־��
func log(i int) {
    logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Key("Key>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Keyf("Key>>>>>>>>>>>>>>>>>>>>>>>>>%v",i)
}