<template>
  <div class="container">
    <div class="about">
      <h1 class="title">О проекте</h1>

      <div class="text">
        <p>
          Наша веб-галерея — это открытое пространство для знакомства с современным и классическим
          изобразительным искусством. Здесь собраны работы художников разных направлений и эпох:
          от реалистических пейзажей до абстрактных композиций. Каждое произведение сопровождается
          информацией об авторе, технике исполнения и времени создания, чтобы дать зрителю
          полный контекст для восприятия.
        </p>
        <p>
          Проект создан для тех, кто хочет открывать для себя новые имена, исследовать
          разнообразие художественных стилей и просто наслаждаться искусством в удобном формате.
          Коллекция постоянно пополняется: мы стремимся к тому, чтобы галерея оставалась
          живым и актуальным местом встречи художников и ценителей прекрасного.
        </p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.container {
  max-width: 1400px;
  margin: 0 auto;
  padding: 0 50px;
}

.about {
  max-width: 660px;
  padding: 80px 0 100px;
}

.title {
  font-size: 2rem;
  font-weight: 600;
  color: #000;
  margin: 0 0 40px;
}

.text {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.text p {
  font-size: 1rem;
  line-height: 1.8;
  color: #444;
}

@media (max-width: 750px) {
  .container { padding: 0 20px; }
  .about { padding: 50px 0 70px; }
  .title { font-size: 1.5rem; margin-bottom: 28px; }
}
</style>
