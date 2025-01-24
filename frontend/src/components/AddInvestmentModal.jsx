import React from "react";
import { Modal, Form, Input, Button } from "antd";

const AddInvestmentModal = ({ visible, onCancel, onAddInvestment }) => {
  const [form] = Form.useForm();

  const handleSubmit = async (values) => {
    console.log(values)
    await onAddInvestment(values);
    form.resetFields();
  };

  return (
    <Modal
      title="Add Investment"
      visible={visible}
      onCancel={onCancel}
      footer={null}
    >
      <Form
        form={form}
        onFinish={handleSubmit}
        layout="vertical"
      >
        <Form.Item
        name="instrument"
        label="Instrument"
        rules={[{ required: true, message: "Please input the instrument name!" }]}
      >
        <Input />
      </Form.Item>

      <Form.Item
        name="qty"
        label="Qty"
        rules={[{ required: true, message: "Please input the quantity!" }]}
      >
        <Input type="number" />
      </Form.Item>

      <Form.Item
        name="avg"
        label="Avg"
        rules={[{ required: true, message: "Please input the price!" }]}
      >
        <Input type="number" />
        </Form.Item>

        <Form.Item>
          <Button type="primary" htmlType="submit" block>
            Add Investment
          </Button>
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default AddInvestmentModal;
