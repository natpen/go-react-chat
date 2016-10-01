import React from 'react';
import ReactDOM from 'react-dom';
import ee from 'event-emitter';
import moment from 'moment';
import uuid from 'uuid';

import LoginForm from './loginForm';
import MessageForm from './messageForm';
import Message from './message';

var emitter = ee({});
var listener;
var ws;

var ChatBox = React.createClass({
  getInitialState: function() {
    return {
      name: '',
      lastActive: null,
      messages: []
    };
  },
  componentDidMount: function() {
    emitter.on('new-messages-emitter', this.handleNewMessages);
    emitter.on('socket-closed-emitter', this.handleSocketClosed);
  },
  componentWillUnmount: function() {
    // ws.send('CLIENT EXIT');
    emitter.off('new-messages-emitter', this.handleNewMessages);
    emitter.off('socket-closed-emitter', this.handleSocketClosed);
  },
  handleNewMessages: function(messages) {
    var processedMessages = JSON.parse(messages.data).map(function(message) {
      return {
        id: message.Id,
        type: message.Type,
        timestamp: message.Timestamp,
        name: message.Name,
        text: message.Text
      };
    });

    // Filter by client-handshake to handle initial lastActive value
    processedMessages
      .filter(function(message) {
        return message.type === 'client-handshake'
      })
      .forEach(function(message) {
        this.setState({
          lastActive: moment(message.timestamp)
        });
      }.bind(this));

    // Filter by user-messages and system-messages to get displayable results
    var newStateMessages = this.state.messages
      .concat(processedMessages)
      .filter(function(message) {
        return message.type === 'user-message' || message.type === 'system-message';
      });
    this.setState({
      messages: newStateMessages
    });
  },
  handleSocketClosed: function(event) {
    this.setState({
      name: '',
      lastActive: null,
      messages: []
    });
  },
  handleNameSubmit: function(name) {
    ws = new WebSocket('ws://' + window.location.host + '/ws'),
    ws.onopen = function open() {
      ws.send(name);
      ws.onmessage = function(data, flags) {
        emitter.emit('new-messages-emitter', data);
      };
      ws.onclose = function(event) {
        emitter.emit('socket-closed-emitter', event);
      }
      this.setState({
        name: name
      });
    }.bind(this);
  },
  handleMessageFormActivity: function() {
    this.setState({
      lastActive: moment(Date.now())
    });
  },
  handleMessageSubmit: function(text) {
    var messageTimestamp = Date.now();

    var message = {
      id: uuid.v4(),
      type: 'user-message',
      timestamp: messageTimestamp,
      name: this.state.name,
      text: text
    };

    ws.send(JSON.stringify(message));

    this.setState({
      lastActive: messageTimestamp
    });
  },
  render: function() {
    if (this.state.name === '') {
      return (
        <div>
          <h1>Go-React-Chat - Login</h1>
          <LoginForm onNameSubmit={this.handleNameSubmit} />
          </div>
        );
    }
    return (
      <div className="chatBox">
        <h1>Go-React-Chat ({this.state.name})</h1>
        <MessageList messages={this.state.messages} lastActive={this.state.lastActive} />
        <MessageForm onMessageFormActivity={this.handleMessageFormActivity} onMessageSubmit={this.handleMessageSubmit} />
        </div>
      );
  }
});

var MessageList = React.createClass({
  render: function() {
    var messageNodes = this.props.messages.map(function(message) {
      return (
        <Message key={message.name + message.timestamp} message={message} lastActive={this.props.lastActive}>
        </Message>
        );
    }.bind(this));
    return (
      <div className="messageList">
        {messageNodes}
      </div>
      );
  }
});

ReactDOM.render(
  <ChatBox />,
  document.getElementById('content')
);
