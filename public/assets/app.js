var app = new Vue({
  el: '#app',
  data: {
    ws: null,
    newMessage: '',
    chatContent: '',
    email: null,
    username: null,
    joined: false
  },

  created: function() {
    var self = this;
    this.ws = new WebSocket('ws://' + window.location.host + '/ws');
    this.ws.addEventListener('message', function(e) {
      var msg = JSON.parse(e.data);
      self.chatContent += '<div class="chip">'
        + '<img src="' + self.gravatarURL(msg.email) + '">'
        + msg.username
      + '</div>'
      + emojione.toImage(msg.message) + '<br/>';

      var element = document.getElementsByClassName("msg_card_body");
      element.scrollTop = element.scrollHeight;
    });
  },

  methods: {
    send: function() {
      if(this.newMessage != '') {
        this.ws.send(
          JSON.stringify({
            email: this.email, 
            username: this.username,
            message: $('<p>').html(this.newMessage).text()
          }
        ));
        this.newMessage = '';
      }
    },

    join: function() {
      if(!this.email) {
        Materialize.toast('Enter an email address', 2000);
        return
      }
      if(!this.username) {
        Materialize.toast('Enter a username', 2000);
        return
      }
    
      this.email = $('<p>').html(this.email).text();
      this.username = $('<p>').html(this.username).text();
      this.joined = true;
    },

    gravatarURL: function(email) {
      return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email);
    }
  }
});
