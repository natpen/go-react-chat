var React      = require('react'),
    ReactDOM   = require('react-dom'),
    ee         = require('event-emitter'),
    emitter    = ee({}), listener,
    moment     = require('moment'),
    Remarkable = require('remarkable'),
    uuid       = require('uuid');

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
    this.setState({ messages: newStateMessages });
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
      this.setState({ name: name });
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
          <h1>Taut - Login</h1>
          <LoginForm onNameSubmit={this.handleNameSubmit} />
          </div>
      );
    }
    return (
        <div className="chatBox">
        <h1>Taut ({this.state.name})</h1>
        <MessageList messages={this.state.messages} lastActive={this.state.lastActive} />
        <MessageForm onMessageFormActivity={this.handleMessageFormActivity} onMessageSubmit={this.handleMessageSubmit} />
        </div>
    );
  }
});

var LoginForm = React.createClass({
  getInitialState: function() {
    return {text: ''};
  },
  handleTextChange: function(e) {
    this.setState({text: e.target.value});
  },
  handleSubmit: function(e) {
    e.preventDefault();
    var text = this.state.text.trim();
    if (!text) {
      return;
    }
    this.props.onNameSubmit(text);
    this.setState({text: ''});
  },
  render: function() {
    return (
        <form className="loginForm" onSubmit={this.handleSubmit}>
        <input
      type="text"
      placeholder="Choose a name..."
      value={this.state.text}
      onChange={this.handleTextChange}
        />
        <input type="submit" value="Log in!" />
        </form>
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

var Message = React.createClass({
  composeMessageLine: function() {
    var message = this.props.message;
    if (message.name.length === 0) return '*' + message.text + '*';

    return '**[**' + moment(message.timestamp).format('HH:mm') + '**]** **' + message.name + '**: ' + message.text;
  },
  rawMarkup: function() {
    var md = new Remarkable();
    var rawMarkup = md.render(this.composeMessageLine());
    return { __html: rawMarkup };
  },
  render: function() {
    // If any non-displayable type messages snuck through, abort the render
    if (this.props.message.type !== 'user-message' && this.props.message.type !== 'system-message') { return null; }

    var className = 'message';
    className += ' ' + this.props.message.type;
    if (this.props.lastActive && moment(this.props.message.timestamp).isAfter(this.props.lastActive)) {
      className += ' new';
    }
    return (
        <div className={className}>
        <span dangerouslySetInnerHTML={this.rawMarkup()} />
        </div>
    );
  }
});

var MessageForm = React.createClass({
  getInitialState: function() {
    return {text: ''};
  },
  handleTextChange: function(e) {
    this.setState({text: e.target.value});
    this.props.onMessageFormActivity();
  },
  handleClick: function(e) {
    this.props.onMessageFormActivity();
  },
  handleSubmit: function(e) {
    e.preventDefault();
    var text = this.state.text.trim();
    if (!text) {
      return;
    }
    this.props.onMessageSubmit(text);
    this.setState({text: ''});
  },
  render: function() {
    return (
        <form className="messageForm" onClick={this.handleClick} onSubmit={this.handleSubmit}>
        <input
      type="text"
      placeholder="Say something..."
      value={this.state.text}
      onChange={this.handleTextChange}
        />
        <input type="submit" value="Post" />
        </form>
    );
  }
});

ReactDOM.render(
    <ChatBox />,
  document.getElementById('content')
);
