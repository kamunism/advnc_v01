package main

import (
 "bufio"
 "fmt"
 "math/rand"
 "os"
 "strings"
 "time"
)

// константы для размера игрового поля и начальных параметров
const (
 boardSize     = 10  // Размер поля 10x10
 initialHealth = 100 // Начальное здоровье
 initialMoves  = 50  // Начальное количество ходов
 numResources  = 8   // Количество ресурсов на поле
 numHazards    = 8   // Количество опасностей на поле
)

// ктруктура игрока
type Player struct {
 X      int // Координата X
 Y      int // Координата Y
 Health int // Здоровье
 Moves  int // Количество ходов
}

// структура игры
type Game struct {
 Board  [][]rune // Двумерный массив, представляющий игровое поле
 Player Player  // Объект игрока
 ExitX  int  // Координата X выхода
 ExitY  int  // Координата Y выхода
}

// initializeBoard инициализирует игровое поле, размещает игрока, выход, ресурсы и опасности
func (g *Game) initializeBoard() {
 // Создаем пустое поле
 g.Board = make([][]rune, boardSize)
 for i := range g.Board {
  g.Board[i] = make([]rune, boardSize)
  for j := range g.Board[i] {
   g.Board[i][j] = ' ' // Заполняем пустые клетки
  }
 }

 // инициализируем генератор случайных чисел
 rand.Seed(time.Now().UnixNano())

 // размещаем игрока в рандомном месте
 g.Player = Player{
  Health: initialHealth,
  Moves:  initialMoves,
 }
 g.Player.X, g.Player.Y = rand.Intn(boardSize), rand.Intn(boardSize)
 // важно: на самом поле клетка игрока не меняется, 'X' отрисовывается сверху в printBoard

 // размещаем выход
 for {
  x, y := rand.Intn(boardSize), rand.Intn(boardSize)
  if g.Board[x][y] == ' ' && (x != g.Player.X || y != g.Player.Y) { // Убеждаемся, что клетка пуста и не занята игроком
   g.Board[x][y] = 'E'
   g.ExitX, g.ExitY = x, y
   break
  }
 }

 // размещаем ресурсы
 for i := 0; i < numResources; i++ {
  for {
   x, y := rand.Intn(boardSize), rand.Intn(boardSize)
   if g.Board[x][y] == ' ' && (x != g.Player.X || y != g.Player.Y) {
    g.Board[x][y] = 'R'
    break
   }
  }
 }

 // растоновка опасности
 for i := 0; i < numHazards; i++ {
  for {
   x, y := rand.Intn(boardSize), rand.Intn(boardSize)
   if g.Board[x][y] == ' ' && (x != g.Player.X || y != g.Player.Y) {
    g.Board[x][y] = 'H'
    break
   }
  }
 }
}

// printBoard выводит текущее состояние игрового поля и статистику игрока
func (g *Game) printBoard() {
 fmt.Println("--- Экспедиция на Чужую Планету ---")
 fmt.Printf("Здоровье: %d | Ходы: %d\n", g.Player.Health, g.Player.Moves)
 fmt.Println("---------------------------------")
 for i := 0; i < boardSize; i++ {
  for j := 0; j < boardSize; j++ {
   if i == g.Player.X && j == g.Player.Y {
    fmt.Printf(" X ") // Отображаем игрока на его текущей позиции
   } else {
    fmt.Printf(" %c ", g.Board[i][j])
   }
  }
  fmt.Println()
 }
 fmt.Println("---------------------------------")
 fmt.Println("Легенда: X=Вы, R=Ресурс (+5 Ходов), H=Опасность (-10 Здоровья), E=Выход")
 fmt.Println("Команды: w (вверх), s (вниз), a (влево), d (вправо), q (выход)")
}

// movePlayer перемещает игрока и обрабатывает взаимодействие с клеткой
func (g *Game) movePlayer(dx, dy int) string {
 newX, newY := g.Player.X+dx, g.Player.Y+dy

 //Проверяем границы поля
 if newX < 0 || newX >= boardSize || newY < 0 || newY >= boardSize {
  return "Заблокировано! Вы не можете выйти за пределы карты."
 }

 g.Player.Moves-- // Каждый ход расходует 1 очко хода

 // получаем содержимое новой клетки
 cellContent := g.Board[newX][newY]

 // обнова нашей позиции
 g.Player.X, g.Player.Y = newX, newY

 // обрабатываем взаимодействие с содержимым клетки
 switch cellContent {
 case 'R':
  g.Player.Moves += 5
  g.Board[newX][newY] = ' ' // Удаляем ресурс после сбора
  return "Вы нашли ресурс! +5 Ходов."
 case 'H':
  g.Player.Health -= 10
  g.Board[newX][newY] = ' ' // Опасность исчезает после взаимодействия
  return "О нет! Опасность! -10 Здоровья."
 case 'E':
  return "Вы нашли Выход!"
 default: // ' ' или другая пустая клетка
  return "Переместились в пустое пространство."
 }
}

// чек закончена ли игра иль нет
func (g *Game) checkGameOver() (bool, string) {
 if g.Player.Health <= 0 {
  return true, "Игра Окончена! У вас закончилось здоровье."
 }
 if g.Player.Moves <= 0 {
  return true, "Игра Окончена! У вас закончились ходы."
 }
 if g.Player.X == g.ExitX && g.Player.Y == g.ExitY {
  return true, "Поздравляем! Вы нашли выход и сбежали с чужой планеты!"
 }
 return false, ""
}

func main() {
 game := Game{}
 game.initializeBoard() // Инициализируем игру

 reader := bufio.NewReader(os.Stdin) // Для чтения ввода пользователя

 for {
  game.printBoard() // Выводим текущее состояние поля
  gameOver, message := game.checkGameOver() // Проверяем условия завершения игры
  if gameOver {
   fmt.Println(message)
   break // Выходим из цикла, если игра окончена
  }

  fmt.Print("Введите ваш ход (w/s/a/d/q): ")
  input, _ := reader.ReadString('\n')
  input = strings.TrimSpace(strings.ToLower(input)) // Очищаем ввод и приводим к нижнему регистру

  var moveMessage string
  switch input {
  case "w":
   moveMessage = game.movePlayer(-1, 0) // Вверх (изменяем X)
  case "s":
   moveMessage = game.movePlayer(1, 0) // Вниз (изменяем X)
  case "a":
   moveMessage = game.movePlayer(0, -1) // Влево (изменяем Y)
  case "d":
   moveMessage = game.movePlayer(0, 1) // Вправо (изменяем Y)
  case "q":
   fmt.Println("Выход из игры. До свидания!")
   return // Выходим из программы
  default:
   moveMessage = "Неверная команда. Используйте w, s, a, d или q."
  }
  fmt.Println(moveMessage)
  fmt.Println() // Пустая строка для лучшей читаемости
 }
}

// компиляция для самых маленьких
// go build (название файла)+.go
// над мини проектом работает Dr0n
