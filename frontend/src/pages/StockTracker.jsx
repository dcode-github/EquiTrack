import React, { useEffect, useState } from "react";
import { useWebSocket } from "../WebSocketContext.js";
import { Table, Typography, Row, Col, Tooltip, Modal, Form, Input, Button, notification } from "antd";
import AddInvestmentModal from "../components/AddInvestmentModal";
import ChildTable from "./ChildTable.jsx";
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
  const API_URL = process.env.REACT_APP_API_URL;
  const WEBSOCKET_URL = process.env.REACT_APP_WEBSOCKET_URL;
  const PORT = process.env.REACT_APP_PORT;

  const [data, setData] = useState([]);
  const [totalInvestmentData, setTotalInvestmentData] = useState({
    total_investment: 0,
    total_currVal: 0,
    total_pnl: 0,
    total_pnl_percent: 0,
  });
  const [loading, setLoading] = useState(true);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [form] = Form.useForm();
  const [pageSize, setPageSize] = useState(10);
  const [childDataMap, setChildDataMap] = useState({});
  const [ws, setWs] = useState(null);
  const [liveDataMap, setLiveDataMap] = useState({});
  const { connect, wsRef } = useWebSocket();

  useEffect(() => {
    fetchInvestments();

    return () => {
      if (ws) ws.close();
    };
  }, []);

  const connectWebSocket = (instrumentList, staticMap) => {
    const instrumentsParam = instrumentList.join(",");
    const socketUrl = `${WEBSOCKET_URL}:${PORT}/priceWebSocket?instrument=${encodeURIComponent(instrumentsParam)}`;

    const socket = connect(socketUrl);

    socket.onopen = () => {
      console.log("WebSocket connected");
    };

    socket.onmessage = (event) => {
      const liveData = JSON.parse(event.data);

      setLiveDataMap((prevLiveDataMap) => ({
        ...prevLiveDataMap,
        [liveData.instrument]: liveData,
      }));
    };

    socket.onclose = () => {
      console.log("WebSocket closed");
    };

    socket.onerror = (error) => {
      console.error("WebSocket error:", error);
    };

    setWs(socket);
  };

  const fetchInvestments = async () => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");

      const response = await fetch(`${API_URL}:${PORT}/investments?userId=${userId}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      const investments = await response.json();
      const staticData = investments.investments;

      const staticMap = {};
      staticData.forEach(item => {
        staticMap[item.instrument] = item;
      });

      setData(staticData);
      setTotalInvestmentData(investments.total_investment_data || {});
      setLoading(false);

      if (staticData.length > 0) {
        const instruments = staticData.map(item => item.instrument);
        connectWebSocket(instruments, staticMap);
      }
    } catch (error) {
      console.error("Error fetching investments:", error);
      setLoading(false);
    }
  };

  const fetchIndividualInvestments = async (instrument) => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");

      const response = await fetch(`${API_URL}:${PORT}/individualInvestments?userId=${userId}&instrument=${instrument}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
      });

      const result = await response.json();
      setChildDataMap(prev => ({
        ...prev,
        [instrument]: result,
      }));
    } catch (error) {
      console.error("Error fetching individual investments:", error);
    }
  };

  const expandedRowRender = (record) => {
    const instrument = record.instrument;

    if (!childDataMap[instrument]) {
      fetchIndividualInvestments(instrument);
    }

    const rowData = childDataMap[instrument];

    return (
      <div>
        {rowData && rowData.length > 0 ? (
          <ChildTable data={rowData} />
        ) : (
          <Text>No data available</Text>
        )}
      </div>
    );
  };

  const addInvestment = async (values) => {
    try {
      const token = sessionStorage.getItem("token");
      const userId = sessionStorage.getItem("userId");
      const qty = parseInt(values.qty, 10);
      const avg = parseFloat(values.avg);
      const user_id = parseInt(userId);

      if (isNaN(qty) || isNaN(avg)) {
        console.error("Invalid qty/avg values.");
        return;
      }

      const response = await fetch(`${API_URL}:${PORT}/investments`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          user_id,
          instrument: values.instrument,
          qty,
          avg,
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
          description: result?.error || "Something went wrong.",
        });
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
          {record.pledged && (
            <Tooltip title="Pledged">
              <Text className="pledged-tag">P: {record.pledged}</Text>
            </Tooltip>
          )}
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
      render: (text) => (
        <Text style={{ color: text > 0 ? "green" : "red" }}>
          {formatNumber(Math.abs(text))}
        </Text>
      ),
    },
    {
      title: "Net chg.",
      dataIndex: "netChng",
      key: "netChng",
      sorter: (a, b) => parseFloat(a.netChng) - parseFloat(b.netChng),
      render: (text) => (
        <Text style={{ color: text > 0 ? "green" : "red" }}>
          {Math.abs(text)}%
        </Text>
      ),
    },
    {
      title: "Day chg.",
      dataIndex: "dayChng",
      key: "dayChng",
      sorter: (a, b) => parseFloat(a.dayChng) - parseFloat(b.dayChng),
      render: (text) => (
        <Text style={{ color: text > 0 ? "green" : "red" }}>
          {Math.abs(text)}%
        </Text>
      ),
    },
  ];

  const updatedData = data.map((staticItem) => {
    const liveItem = liveDataMap[staticItem.instrument] || {};

    return {
      ...staticItem,
      ltp: liveItem.price || 0,
      dayChng: liveItem.per_change || 0,
      currVal: (liveItem.price || 0) * staticItem.qty,
      pnl: (liveItem.price || 0) * staticItem.qty - staticItem.tot_invest,
      netChng: liveItem.price
        ? (((liveItem.price * staticItem.qty - staticItem.tot_invest) / staticItem.tot_invest) * 100).toFixed(2)
        : 0,
    };
  });

  useEffect(() => {
    const total_investment = updatedData.reduce((acc, item) => acc + item.tot_invest, 0);
    const total_currVal = updatedData.reduce((acc, item) => acc + item.currVal, 0);
    const total_pnl = total_currVal - total_investment;
    const total_pnl_percent = total_investment !== 0 ? ((total_pnl / total_investment) * 100).toFixed(2) : 0;
  
    setTotalInvestmentData({
      total_investment,
      total_currVal,
      total_pnl,
      total_pnl_percent,
    });
  }, [updatedData]);

  return (
    <div className="stock-tracker-container">
      <Navbar />
      <div className="summary-container">
        <Row gutter={[16, 16]}>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>
                {formatNumber(totalInvestmentData.total_investment || 0)}
              </Title>
              <Text>Total Investment</Text>
            </div>
          </Col>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>
                {formatNumber(totalInvestmentData.total_currVal || 0)}
              </Title>
              <Text>Current Value</Text>
            </div>
          </Col>
          <Col span={10} className="summary-item">
            <div>
              <Title level={5}>
                {formatNumber(totalInvestmentData.total_pnl || 0)}{" "}
                <Text type={totalInvestmentData.total_pnl > 0 ? "success" : "danger"}>
                  {totalInvestmentData.total_pnl_percent > 0
                    ? `+${totalInvestmentData.total_pnl_percent}%`
                    : `${totalInvestmentData.total_pnl_percent}%`}
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
          dataSource={updatedData}
          rowKey={(record) => record.instrument}
          loading={loading}
          pagination={{
            pageSize,
            pageSizeOptions: ["5", "10", "20", "50"],
            showSizeChanger: true,
            onShowSizeChange: (current, size) => setPageSize(size),
          }}
          expandable={{
            expandedRowRender,
            rowExpandable: () => true,
          }}
          className="stock-table"
          sticky
          bordered
        />
      </div>
    </div>
  );
};

export default StockTracker;
