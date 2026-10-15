
package main

import (
	"fmt"
	"strings"
)

// 1. FACADE — ФАСАД


type Lights struct{}

func (l *Lights) On() {
	fmt.Println("Свет включен")
}

type Projector struct{}

func (p *Projector) On() {
	fmt.Println("Проектор включен")
}

type SoundSystem struct{}

func (s *SoundSystem) On() {
	fmt.Println("Звуковая система включена")
}


type HomeTheaterFacade struct {
	lights    *Lights
	projector *Projector
	sound     *SoundSystem
}

func NewHomeTheater() *HomeTheaterFacade {
	return &HomeTheaterFacade{
		lights:    &Lights{},
		projector: &Projector{},
		sound:     &SoundSystem{},
	}
}

func (h *HomeTheaterFacade) WatchMovie() {
	fmt.Println("\n--- 1. Фасад ---")
	h.lights.On()
	h.projector.On()
	h.sound.On()
	fmt.Println("Кинотеатр готов к просмотру")
}


// 2. BUILDER — СТРОИТЕЛЬ


type Computer struct {
	CPU     string
	RAM     int
	Storage int
	GPU     string
}

type ComputerBuilder struct {
	computer Computer
}

func NewComputerBuilder() *ComputerBuilder {
	return &ComputerBuilder{}
}

func (b *ComputerBuilder) SetCPU(cpu string) *ComputerBuilder {
	b.computer.CPU = cpu
	return b
}

func (b *ComputerBuilder) SetRAM(ram int) *ComputerBuilder {
	b.computer.RAM = ram
	return b
}

func (b *ComputerBuilder) SetStorage(storage int) *ComputerBuilder {
	b.computer.Storage = storage
	return b
}

func (b *ComputerBuilder) SetGPU(gpu string) *ComputerBuilder {
	b.computer.GPU = gpu
	return b
}

func (b *ComputerBuilder) Build() Computer {
	return b.computer
}


// 3. VISITOR — ПОСЕТИТЕЛЬ



type Visitor interface {
	VisitBook(*Book)
	VisitMovie(*Movie)
}

type Visitable interface {
	Accept(Visitor)
}

type Book struct {
	Title string
	Price float64
}

func (b *Book) Accept(v Visitor) {
	v.VisitBook(b)
}

type Movie struct {
	Title    string
	Duration int
}

func (m *Movie) Accept(v Visitor) {
	v.VisitMovie(m)
}


type InfoVisitor struct{}

func (v *InfoVisitor) VisitBook(b *Book) {
	fmt.Printf("Книга: %s, цена: %.2f\n", b.Title, b.Price)
}

func (v *InfoVisitor) VisitMovie(m *Movie) {
	fmt.Printf("Фильм: %s, длительность: %d мин.\n",
		m.Title, m.Duration)
}


// 4. COMMAND — КОМАНДА


type Command interface {
	Execute()
}

type TV struct{}

func (t *TV) TurnOn() {
	fmt.Println("Телевизор включен")
}

type TurnOnTVCommand struct {
	tv *TV
}

func (c *TurnOnTVCommand) Execute() {
	c.tv.TurnOn()
}

type RemoteControl struct {
	command Command
}

func (r *RemoteControl) SetCommand(c Command) {
	r.command = c
}

func (r *RemoteControl) PressButton() {
	if r.command != nil {
		r.command.Execute()
	}
}


// 5. CHAIN OF RESPONSIBILITY — ЦЕПОЧКА ВЫЗОВОВ


type Request struct {
	User  string
	Admin bool
}

type Handler interface {
	SetNext(Handler)
	Handle(*Request) bool
}

type BaseHandler struct {
	next Handler
}

func (b *BaseHandler) SetNext(next Handler) {
	b.next = next
}

func (b *BaseHandler) HandleNext(r *Request) bool {
	if b.next != nil {
		return b.next.Handle(r)
	}
	return true
}

type AuthHandler struct {
	BaseHandler
}

func (h *AuthHandler) Handle(r *Request) bool {
	if strings.TrimSpace(r.User) == "" {
		fmt.Println("Ошибка: пользователь не авторизован")
		return false
	}
	fmt.Println("Проверка авторизации пройдена")
	return h.HandleNext(r)
}

type AdminHandler struct {
	BaseHandler
}

func (h *AdminHandler) Handle(r *Request) bool {
	if !r.Admin {
		fmt.Println("Ошибка: нужны права администратора")
		return false
	}
	fmt.Println("Проверка прав пройдена")
	return h.HandleNext(r)
}

type FinalHandler struct{}

func (h *FinalHandler) SetNext(next Handler) {}

func (h *FinalHandler) Handle(r *Request) bool {
	fmt.Println("Запрос успешно обработан")
	return true
}


// 6. FACTORY METHOD — ФАБРИЧНЫЙ МЕТОД


type Notifier interface {
	Send(message string)
}

type EmailNotifier struct{}

func (e *EmailNotifier) Send(message string) {
	fmt.Println("Email:", message)
}

type SMSNotifier struct{}

func (s *SMSNotifier) Send(message string) {
	fmt.Println("SMS:", message)
}

type NotifierFactory interface {
	CreateNotifier() Notifier
}

type EmailFactory struct{}

func (f *EmailFactory) CreateNotifier() Notifier {
	return &EmailNotifier{}
}

type SMSFactory struct{}

func (f *SMSFactory) CreateNotifier() Notifier {
	return &SMSNotifier{}
}

func Notify(factory NotifierFactory, message string) {
	notifier := factory.CreateNotifier()
	notifier.Send(message)
}


// 7. STRATEGY — СТРАТЕГИЯ

type PaymentStrategy interface {
	Pay(amount float64)
}

type CardPayment struct{}

func (c *CardPayment) Pay(amount float64) {
	fmt.Printf("Оплата картой: %.2f\n", amount)
}

type CashPayment struct{}

func (c *CashPayment) Pay(amount float64) {
	fmt.Printf("Оплата наличными: %.2f\n", amount)
}

type PaymentContext struct {
	strategy PaymentStrategy
}

func (p *PaymentContext) SetStrategy(s PaymentStrategy) {
	p.strategy = s
}

func (p *PaymentContext) Pay(amount float64) {
	if p.strategy == nil {
		fmt.Println("Способ оплаты не выбран")
		return
	}
	p.strategy.Pay(amount)
}


// 8. STATE — СОСТОЯНИЕ


type OrderState interface {
	Next(*Order)
	Name() string
}

type NewOrderState struct{}

func (s *NewOrderState) Name() string {
	return "Новый"
}

func (s *NewOrderState) Next(o *Order) {
	fmt.Println("Заказ передан в обработку")
	o.state = &ProcessingState{}
}

type ProcessingState struct{}

func (s *ProcessingState) Name() string {
	return "В обработке"
}

func (s *ProcessingState) Next(o *Order) {
	fmt.Println("Заказ отправлен")
	o.state = &ShippedState{}
}

type ShippedState struct{}

func (s *ShippedState) Name() string {
	return "Отправлен"
}

func (s *ShippedState) Next(o *Order) {
	fmt.Println("Заказ уже отправлен: переход невозможен")
}

type Order struct {
	state OrderState
}

func NewOrder() *Order {
	return &Order{state: &NewOrderState{}}
}

func (o *Order) Next() {
	o.state.Next(o)
}

func (o *Order) ShowState() {
	fmt.Println("Текущее состояние заказа:", o.state.Name())
}




func main() {
	// 1. Фасад
	theater := NewHomeTheater()
	theater.WatchMovie()

	// 2. Строитель
	fmt.Println("\n--- 2. Строитель ---")
	pc := NewComputerBuilder().
		SetCPU("AMD Ryzen 7").
		SetRAM(32).
		SetStorage(1000).
		SetGPU("RTX 5070").
		Build()
	fmt.Printf("Компьютер: %+v\n", pc)

	// 3. Посетитель
	fmt.Println("\n--- 3. Посетитель ---")
	visitor := &InfoVisitor{}
	items := []Visitable{
		&Book{Title: "Go для начинающих", Price: 1500},
		&Movie{Title: "Матрица", Duration: 136},
	}
	for _, item := range items {
		item.Accept(visitor)
	}

	// 4. Команда
	fmt.Println("\n--- 4. Команда ---")
	tv := &TV{}
	remote := &RemoteControl{}
	remote.SetCommand(&TurnOnTVCommand{tv: tv})
	remote.PressButton()

	// 5. Цепочка вызовов
	fmt.Println("\n--- 5. Цепочка вызовов ---")
	auth := &AuthHandler{}
	admin := &AdminHandler{}
	final := &FinalHandler{}
	auth.SetNext(admin)
	admin.SetNext(final)

	request := &Request{User: "slava", Admin: true}
	auth.Handle(request)

	// 6. Фабричный метод
	fmt.Println("\n--- 6. Фабричный метод ---")
	Notify(&EmailFactory{}, "Заказ оформлен")
	Notify(&SMSFactory{}, "Заказ отправлен")

	// 7. Стратегия
	fmt.Println("\n--- 7. Стратегия ---")
	payment := &PaymentContext{}
	payment.SetStrategy(&CardPayment{})
	payment.Pay(5000)
	payment.SetStrategy(&CashPayment{})
	payment.Pay(2500)

	// 8. Состояние
	fmt.Println("\n--- 8. Состояние ---")
	order := NewOrder()
	order.ShowState()
	order.Next()
	order.ShowState()
	order.Next()
	order.ShowState()
}
