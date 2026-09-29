package mqttbridge

type MqttSettings interface {
	Host() string
	User() string
	Password() string
	CaCert() string
}

/*
func loadCert(caCert string) (*x509.Certificate, error) {
	data, err := os.ReadFile(caCert)
	if err != nil {
		return nil, fmt.Errorf("can not read certificate: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("failed to decode pem in %s", caCert)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	return cert, nil
}

func newTLSConfig(caCert string) (*tls.Config, error) {
	cert, err := loadCert(caCert)
	if err != nil {
		return nil, err
	}
	result := tls.Config{}
	result.RootCAs = x509.NewCertPool()
	result.RootCAs.AddCert(cert)
	return &result, nil
}

func plainClientOptions(s MqttSettings, ch mqtt.OnConnectHandler) (mqtt.ClientOptions, error) {
	return mqtt.ClientOptions{
		Servers:        []*url.URL{{Scheme: "tcp", Host: s.Host()}},
		AutoReconnect:  true,
		OnConnect:      ch,
		ClientID:       "thermo-watchdog",
		CleanSession:   true,
		PingTimeout:    10 * time.Second,
		ConnectTimeout: time.Second,
		KeepAlive:      10,
		Username:       s.User(),
		Password:       s.Password(),
	}, nil
}

func secureClientOptions(s MqttSettings, handler mqtt.OnConnectHandler) (mqtt.ClientOptions, error) {
	tc, err := newTLSConfig(s.CaCert())
	if err != nil {
		return mqtt.ClientOptions{}, fmt.Errorf("can not create tls config: %w", err)
	}
	return mqtt.ClientOptions{
		Servers:        []*url.URL{{Scheme: "tcps", Host: s.Host()}},
		AutoReconnect:  true,
		OnConnect:      handler,
		ClientID:       "thermo-watchdog",
		CleanSession:   true,
		PingTimeout:    10 * time.Second,
		ConnectTimeout: time.Second,
		KeepAlive:      10,
		TLSConfig:      tc,
		Username:       s.User(),
		Password:       s.Password(),
	}, nil
}

func newClientOptions(s MqttSettings, handler mqtt.OnConnectHandler) (mqtt.ClientOptions, error) {
	if s.CaCert() == "" {
		return plainClientOptions(s, handler)
	}
	return secureClientOptions(s, handler)
}

func connect(s MqttSettings, handler mqtt.OnConnectHandler) (mqtt.Client, error) {
	opts, err := newClientOptions(s, handler)
	if err != nil {
		return nil, fmt.Errorf("mqtt.connect: %w", err)
	}
	client := mqtt.NewClient(&opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, fmt.Errorf("mqtt.connect: %w", token.Error())
	}
	return client, nil
}

*/
