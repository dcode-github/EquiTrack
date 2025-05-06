// src/components/IndividualInvestmentTable.js
import React from "react";
import { Table, Typography } from "antd";

const { Text } = Typography;

const formatNumber = (num) => {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 2,
  }).format(num);
};

const IndividualInvestmentTable = ({ data }) => {
  // console.log("data is "+data);
  const columns = [
    {
      title: "Instrument",
      dataIndex: "instrument",
      key: "instrument",
    },
    {
      title: "Quantity",
      dataIndex: "qty",
      key: "qty",
    },
    {
      title: "Average Cost",
      dataIndex: "avg",
      key: "avg",
      render: (text) => formatNumber(text),
    },
    {
      title: "Date",
      dataIndex: "date",
      key: "date",
    },
  ];

  return (
    <Table
      columns={columns}
      dataSource={data} // This will be passed as a prop
      rowKey="date"
      pagination={false}
    />
  );
};

export default IndividualInvestmentTable;
