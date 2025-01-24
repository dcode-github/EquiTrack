import React, { useEffect, useState } from "react";
import { Table, Typography, Row, Col, Tooltip, Modal, Form, Input, Button, notification } from "antd";
import AddInvestmentModal from "../components/AddInvestmentModal";
import Navbar from "../components/Navbar";
import "./StockTracker.css";

const { Title, Text } = Typography;

const formatNumber = (num) => {
  return new Intl.NumberFormat("en-IN", {
    style: "currency",
    currency: "INR",
    maximumFractionDigits: 2,
  }).format(num);
};

const StockTracker = () => {
  const [data, setData] = useState([]);
  const [totalInvestmentData, setTotalInvestmentData] = useState({
    totalInvestment: 0,
    totalCurrentVal: 0,
    totalPNL: 0,
    totalPNLPercent: 0,
  });
  const [loading, setLoading] = useState(true);
  const [expandedRowKeys, setExpandedRowKeys] = useState([]);
  const [intervalId, setIntervalId] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [pageSize, setPageSize] = useState(10);

  useEffect(() => {
    fetchInvestments();
  }, []);

  const fetchInvestments = async () => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");

      const response = await fetch(`http://localhost:8080/investments?userId=${userId}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      const investments = await response.json();

      setData(investments.investments);
      setTotalInvestmentData(investments.total_investment_data);
      setLoading(false);
    } catch (error) {
      console.error("Error fetching investments:", error);
      setLoading(false);
    }
  };

  const isWithinTradingHours = () => {
    const currentHour = new Date().getHours();
    return currentHour >= 9 && currentHour < 16;
  };

  useEffect(() => {
    if (isWithinTradingHours()) {
      fetchInvestments();
      const id = setInterval(() => {
        fetchInvestments();
      }, 2000);
      setIntervalId(id);
    }
    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, []);

  useEffect(() => {
    if (!isWithinTradingHours() && intervalId) {
      clearInterval(intervalId);
    }
  }, [intervalId]);

  const fetchIndividualInvestments = async (instrument) => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");

      const response = await fetch(`http://localhost:8080/individualInvestments?userId=${userId}&instrument=${instrument}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });
      const individualInvestments = await response.json();
      return individualInvestments;
    } catch (error) {
      console.error("Error fetching individual investments:", error);
      return [];
    }
  };

  const expandedRowRender = async (record) => {
    const individualInvestments = await fetchIndividualInvestments(record.instrument);
    return (
      <Table
        columns={[
          {
            title: "Date",
            dataIndex: "date",
            key: "date",
          },
          {
            title: "Quantity",
            dataIndex: "qty",
            key: "qty",
            render: (text) => formatNumber(text).replace("₹", ""),
          },
          {
            title: "Price",
            dataIndex: "avg",
            key: "avg",
            render: (text) => formatNumber(text),
          },
        ]}
        dataSource={individualInvestments}
        rowKey={(record) => record.date}
        pagination={false}
      />
    );
  };

  const addInvestment = async (values) => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");
      const qty = parseInt(values.qty, 10);
      const user_id = parseInt(userId);
      const avg = parseFloat(values.avg);
      if (isNaN(qty) || isNaN(avg)) {
        console.error("Invalid data: qty should be an integer and avg should be a float.");
        return;
      }
      const response = await fetch("http://localhost:8080/investments", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          user_id: user_id,
          instrument: values.instrument,
          qty: qty,
          avg: avg,
        }),
      });

      const result = await response.json();

      if (response.ok) {
        fetchInvestments();
        setIsModalVisible(false);
        form.resetFields();
        notification.success({
          message: "Investment Added",
          description: `Your investment in ${values.instrument} has been successfully added.`,
        });
      } else {
        notification.error({
          message: "Error Adding Investment",
          description: `There was an error in recording the investment. Please try again.`,
        });
        console.error("Failed to add investment:", result);
      }
    } catch (error) {
      console.error("Error adding investment:", error);
    }
  };

  const columns = [
    {
      title: "Instrument",
      dataIndex: "instrument",
      key: "instrument",
      sorter: (a, b) => a.instrument.localeCompare(b.instrument),
      render: (text, record) => (
        <div>
          <Text>
            <strong>{text}</strong>
          </Text>
          {record.pledged && <Tooltip title="Pledged"><Text className="pledged-tag">P: {record.pledged}</Text></Tooltip>}
        </div>
      ),
    },
    {
      title: "Qty.",
      dataIndex: "qty",
      key: "qty",
      sorter: (a, b) => a.qty - b.qty,
    },
    {
      title: "Avg. cost",
      dataIndex: "avg",
      key: "avg",
      sorter: (a, b) => a.avg - b.avg,
      render: (text) => formatNumber(text),
    },
    {
      title: "LTP",
      dataIndex: "ltp",
      key: "ltp",
      sorter: (a, b) => a.ltp - b.ltp,
      render: (text) => formatNumber(text),
    },
    {
      title: "Investment",
      dataIndex: "tot_invest",
      key: "tot_invest",
      sorter: (a, b) => a.tot_invest - b.tot_invest,
      render: (text) => formatNumber(text),
    },
    {
      title: "Cur. val",
      dataIndex: "currVal",
      key: "currVal",
      sorter: (a, b) => a.currVal - b.currVal,
      render: (text) => formatNumber(text),
    },
    {
      title: "P&L",
      dataIndex: "pnl",
      key: "pnl",
      sorter: (a, b) => a.pnl - b.pnl,
      render: (text) => <Text style={{ color: text > 0 ? "green" : "red" }}>{formatNumber(Math.abs(text))}</Text>,
    },
    {
      title: "Net chg.",
      dataIndex: "netChng",
      key: "netChng",
      sorter: (a, b) => a.netChng - b.netChng,
      render: (text) => <Text style={{ color: text > 0 ? "green" : "red" }}>{Math.abs(text)}%</Text>,
    },
    {
      title: "Day chg.",
      dataIndex: "dayChng",
      key: "dayChng",
      sorter: (a, b) => a.dayChng - b.dayChng,
      render: (text) => <Text style={{ color: text > 0 ? "green" : "red" }}>{Math.abs(text)}%</Text>,
    },
  ];

  return (
    <div className="stock-tracker-container">
      <Navbar />
      <div className="summary-container">
        <Row gutter={[16, 16]}>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>{formatNumber(totalInvestmentData.total_investment)}</Title>
              <Text>Total Investment</Text>
            </div>
          </Col>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>{formatNumber(totalInvestmentData.total_currVal)}</Title>
              <Text>Current Value</Text>
            </div>
          </Col>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>
                {formatNumber(totalInvestmentData.total_pnl)}{" "}
                <Text type={totalInvestmentData.total_pnl > 0 ? "success" : "danger"}>
                  ({totalInvestmentData.total_pnl_percent > 0 ? `+${totalInvestmentData.total_pnl_percent}%` : `${totalInvestmentData.total_pnl_percent}%`})
                </Text>
              </Title>
              <Text>P&L</Text>
            </div>
          </Col>
          <Button
            type="primary"
            onClick={() => setIsModalVisible(true)}
            className="add-investment-button"
          >
            + Investment
          </Button>
        </Row>
      </div>
      <AddInvestmentModal
        visible={isModalVisible}
        onCancel={() => setIsModalVisible(false)}
        onAddInvestment={addInvestment}
      />

      <div className="table-container">
        <Table
          columns={columns}
          dataSource={data}
          rowKey={(record) => record.instrument}
          loading={loading}
          pagination={{
            pageSize: pageSize,
            pageSizeOptions: ['5', '10', '20', '50'],
            showSizeChanger: true,
            onShowSizeChange: (current, size) => setPageSize(size),
          }}
          expandable={{
            expandedRowRender,
            rowExpandable: (record) => record.expandable !== false,
          }}
          className="stock-table"
          sticky={true}
          bordered
        />
      </div>
    </div>
  );
};

export default StockTracker;
