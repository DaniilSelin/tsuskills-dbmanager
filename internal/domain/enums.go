type EmploymentType string

const (
  ПолнаяЗанятость      EmploymentType = "Полная занятость"
  ЧастичнаяЗанятость   EmploymentType = "Частичная занятость"
  ВременнаяЗанятость   EmploymentType = "Временная занятость"
  Стажировка           EmploymentType = "Стажировка"
  РазоваяПомощь        EmploymentType = "Разовая помощь"
)

// switch при компиляции создать jump table
func (e EmploymentType) IsValid() bool {
  switch e {
  case ПолнаяЗанятость, ЧастичнаяЗанятость, ВременнаяЗанятость,
       Стажировка, РазоваяПомощь:
    return true
  }
  return false
}

type WorkSchedule string

const (
  ОфисныйГрафик    WorkSchedule = "Офисный"
  Удалённый        WorkSchedule = "Удалённый"
  СменныйГрафик    WorkSchedule = "Сменный график"
  ГибкийГрафик     WorkSchedule = "Гибкий график"
)

func (w WorkSchedule) IsValid() bool {
  switch w {
  case ОфисныйГрафик, Удалённый, СменныйГрафик, ГибкийГрафик:
    return true
  }
  return false
}

type CompensationType string

const (
  Зарплата         CompensationType = "Зарплата"
  Стипендия        CompensationType = "Стипендия"
  Премия           CompensationType = "Премия"
  Баллы            CompensationType = "Баллы"
  БезВознаграждения CompensationType = "Без вознаграждения"
)

func (c CompensationType) IsValid() bool {
  switch c {
  case Зарплата, Стипендия, Премия, Баллы, БезВознаграждения:
    return true
  }
  return false
}