package memoria

import (
	"testing"
	"time"
)

func TestTiempoDeVidaMáximo(t *testing.T) {
	var almacén = NuevoAlmacén(3*time.Second, 1*time.Second)
	almacén.Almacenar("1", "objeto1")

	var contador = 0
	for {
		time.Sleep(500 * time.Millisecond)
		contador++

		obj1, existe1 := almacén.Adquirir("1")
		if contador == 6 {
			// Transcurridos más de tres segundos
			noDebeExistir(t, existe1, obj1)
			break
		} else {
			debeExistir(t, existe1, obj1, "objeto1")
		}
	}
}

func TestTiempoDeVidaSinUso(t *testing.T) {
	var almacén = NuevoAlmacén(3*time.Second, 1*time.Second)
	almacén.Almacenar("1", "objeto1")
	almacén.Almacenar("2", "objeto2")

	time.Sleep(800 * time.Millisecond)

	obj1, existe1 := almacén.Adquirir("1")
	debeExistir(t, existe1, obj1, "objeto1")

	// Esperar a que expire el «2».
	// El «1» no debe expirar aún.
	time.Sleep(300 * time.Millisecond)

	obj1, existe1 = almacén.Adquirir("1")
	debeExistir(t, existe1, obj1, "objeto1")

	obj2, existe2 := almacén.Adquirir("2")
	noDebeExistir(t, existe2, obj2)

	// Esperar más de un segundo para que expire el restante
	time.Sleep(1100 * time.Millisecond)

	obj1, existe1 = almacén.Adquirir("1")
	noDebeExistir(t, existe1, obj1)
}

func TestSobrescritura(t *testing.T) {
	var almacén = NuevoAlmacén(3*time.Second, 1*time.Second)
	almacén.Almacenar("I", "valor1")

	time.Sleep(800 * time.Millisecond)

	obj, existe := almacén.Adquirir("I")
	debeExistir(t, existe, obj, "valor1")

	time.Sleep(800 * time.Millisecond)

	almacén.Almacenar("I", "valor2")

	time.Sleep(800 * time.Millisecond)

	obj, existe = almacén.Adquirir("I")
	debeExistir(t, existe, obj, "valor2")

	time.Sleep(1100 * time.Millisecond)

	obj, existe = almacén.Adquirir("I")
	noDebeExistir(t, existe, obj)
}

func debeExistir(t *testing.T, existe bool, objeto interface{}, valorEsperado string) {
	if !existe {
		t.Error("objeto no encontrado")
	}
	if objeto.(string) != valorEsperado {
		t.Error("valor de objeto inesperado")
	}
}

func noDebeExistir(t *testing.T, existe bool, objeto interface{}) {
	if existe {
		t.Error("objeto no expirado")
	}
	if objeto != nil {
		t.Error("valor de objeto expirado no nulo")
	}
}
