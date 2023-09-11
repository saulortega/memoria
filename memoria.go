package memoria

import (
	"sync"
	"time"
)

// Almacén representa un conjunto de objetos con un tiempo de vida común.
type Almacén struct {
	mu        sync.Mutex
	tdvm      time.Duration // Tiempo de vida máximo de los objetos almacenados.
	tdvsu     time.Duration // Tiempo de vida de los objetos almacenados desde la última vez que fueron adquiridos.
	elementos map[string]*elemento
}

type elemento struct {
	tdvm   *time.Timer
	tdvsu  *time.Timer
	objeto interface{}
}

// NuevoAlmacén crea un nuevo Almacén donde se almacenarán los objetos.
// Los objetos almacenados se eliminarán después de transcurrir tiempoDeVidaSinUso desde el uso más reciente,
// o, al transcurrir tiempoDeVidaMáximo desde el momento en que se almacenó, lo primero que ocurra.
func NuevoAlmacén(tiempoDeVidaMáximo time.Duration, tiempoDeVidaSinUso time.Duration) *Almacén {
	return &Almacén{
		tdvm:      tiempoDeVidaMáximo,
		tdvsu:     tiempoDeVidaSinUso,
		elementos: map[string]*elemento{},
	}
}

// Almacenar almacena un objeto asignándole un identificador único.
// Si ya existe un objeto con el mismo identificador, se sobrescribirá.
func (O *Almacén) Almacenar(identificador string, objeto interface{}) {
	O.mu.Lock()

	var elementoExistente, existe = O.elementos[identificador]
	if existe {
		reiniciarTemporizador(elementoExistente.tdvm, O.tdvm)
		reiniciarTemporizador(elementoExistente.tdvsu, O.tdvsu)
		elementoExistente.objeto = objeto
		O.elementos[identificador] = elementoExistente
		O.mu.Unlock()
		return
	}

	var elementoNuevo = &elemento{
		tdvm:   time.NewTimer(O.tdvm),
		tdvsu:  time.NewTimer(O.tdvsu),
		objeto: objeto,
	}

	go func() {
		<-elementoNuevo.tdvm.C
		O.removerElemento(identificador, elementoNuevo.tdvsu)
	}()

	go func() {
		<-elementoNuevo.tdvsu.C
		O.removerElemento(identificador, elementoNuevo.tdvm)
	}()

	O.elementos[identificador] = elementoNuevo

	O.mu.Unlock()
}

// Adquirir obtiene un objeto dado su identificador.
// El valor lógico retornado indica si existe o no.
func (O *Almacén) Adquirir(identificador string) (interface{}, bool) {
	O.mu.Lock()
	defer O.mu.Unlock()

	var elemento, existe = O.elementos[identificador]
	if !existe {
		return nil, false
	}

	reiniciarTemporizador(elemento.tdvsu, O.tdvsu)

	return elemento.objeto, true
}

func (O *Almacén) removerElemento(identificador string, temporizadorADetener *time.Timer) {
	O.mu.Lock()

	if !temporizadorADetener.Stop() {
		<-temporizadorADetener.C
	}

	delete(O.elementos, identificador)

	O.mu.Unlock()
}

func reiniciarTemporizador(temporizador *time.Timer, duración time.Duration) {
	if !temporizador.Stop() {
		<-temporizador.C
	}

	temporizador.Reset(duración)
}
