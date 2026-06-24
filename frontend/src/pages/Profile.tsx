import ChangePassword from "../components/ChangePassword"
import UpdateUserForm from "../components/UpdateUserForm"

function Profile() {

  return (
    <div className="max-w-xl mx-auto py-8 px-4 flex flex-col gap-6">
        <UpdateUserForm />
        <ChangePassword />
    </div>
  )
}
export default Profile