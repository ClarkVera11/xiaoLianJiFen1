const formatPhoneNumber = (phone: string) => {
  if (!phone) return '';
  return phone.replace(/(\d{3})\d{4}(\d{4})/, '$1****$2');
};

// 在表格列定义中修改 phone 的渲染方式
{
  title: '手机号',
  dataIndex: 'phone',
  key: 'phone',
  render: (text: string) => formatPhoneNumber(text)
}, 